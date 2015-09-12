package hikarianicmp

import (
	"encoding/binary"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
	MagicCode        = 55
	AckCode          = 56
)

type HikarianIcmp struct {
	client, server string
	mode           string
	TCPPool        *TCPConnPool
	ChannelPool    *ChannelPool
}

func NewHikarianIcmp(client, server, mode string) *HikarianIcmp {
	return &HikarianIcmp{
		server:      server,
		client:      client,
		mode:        mode,
		TCPPool:     NewTCPConnPool(),
		ChannelPool: NewChannelPool(),
	}
}

func (self *HikarianIcmp) transportServer(clientConn *icmp.PacketConn, caddr net.Addr, icmpChannel chan []byte) {
	for {
		body, ok := <-icmpChannel
		if ok == false {
			return
		}
		hash := binary.BigEndian.Uint16(body[0:2]) + binary.BigEndian.Uint16(body[2:4])
		serverConn := self.TCPPool.Get(hash)
		if serverConn == nil {
			server, err := net.ResolveTCPAddr("tcp", self.server)
			if err != nil {
				log.Fatalln("resolve server address error: ", err.Error())
			}
			serverConn, err = net.DialTCP("tcp", nil, server)
			if err != nil {
				log.Println("connect remote address error:", err)
				return
			}
			// serverConn.SetDeadline(time.Now().Add(time.Second * 5))
			self.TCPPool.Set(hash, serverConn)
		}
		nw, err := serverConn.Write(body[4:])
		if err != nil {
			log.Println("write server error: ", err.Error())
			return
		}
		log.Println("get echo reply size ", nw)
		readChannel := make(chan []byte)
		go func() {
			rb := make([]byte, 1024)
			for {
				nr, err := serverConn.Read(rb)
				if err != nil && err != io.EOF {
					log.Println("read server error: ", err.Error())
					close(readChannel)
					return
				}
				log.Println("read ", nr)
				readChannel <- rb[:nr]
			}
		}()
		go func() {
			for {
				wb, ok := <-readChannel
				if ok == false {
					return
				}
				log.Println("read from channel ", len(wb))
				reply, err := (&icmp.Message{
					Type: ipv4.ICMPTypeEchoReply,
					Code: MagicCode,
					Body: &icmp.Echo{
						ID:   int(binary.BigEndian.Uint16((body[0:2]))),
						Seq:  int(binary.BigEndian.Uint16((body[2:4]))),
						Data: wb,
					},
				}).Marshal(nil)
				if err != nil {
					log.Println("marshal echo reply error: ", err.Error())
					return
				}
				ReSend:
				for i := 0; i < 3; i++ {
					numWrite, err := clientConn.WriteTo(reply, caddr)
					if err != nil {
						log.Println("write echo reply error: ", err.Error())
						return
					}
					log.Println("write echo reply size ", numWrite)

					select {
					case _ = <-icmpChannel:
						log.Println("read ack")
						break ReSend
					case <-time.After(2 * time.Second):
						log.Println("timeout")
						continue
					}
				}
				log.Println("break")
			}
		}()
	}
}

func (self *HikarianIcmp) transportClient(clientConn *net.TCPConn) {
	rb := make([]byte, 10240)
	wb := make([]byte, 10240)
	body := make([]byte, 10240)
	host, port, err := net.SplitHostPort(clientConn.RemoteAddr().String())
	if err != nil {
		log.Fatal("split host port error: ", err.Error())
		return
	}
	ip4 := strings.Split(host, ".")
	id := 0
	for index, value := range ip4 {
		log.Println(index, value)
		i, err := strconv.Atoi(value)
		if err != nil {
			log.Println("strconv ip error: ", err.Error())
			continue
		}
		id += int(uint(i) >> uint(index))
	}
	seq, err := strconv.Atoi(port)
	if err != nil {
		log.Println("strconv port error: ", err.Error())
	}

	for {
		size := 0
		nr, err := clientConn.Read(rb)
		if nr == 0 {
			return
		}
		log.Println("get client data ", rb[:nr])
		if nr > 0 {
			wb = append(wb[:size], rb[:nr]...)
			size += nr
		}
		if err != nil && err != io.EOF {
			log.Println("read server error: ", err.Error())
			return
		}
		laddr := &net.IPAddr{IP: net.ParseIP("0.0.0.0")}
		raddr, err := net.ResolveIPAddr("ip", self.server)
		if err != nil {
			log.Fatalln("parse remote addr error: ", err.Error())
		}

		serverConn, err := net.DialIP("ip4:icmp", laddr, raddr)
		if err != nil {
			log.Fatalln("dial ip failed", err.Error())
			return
		}
		defer serverConn.Close()

		log.Printf("set echo request id:%d and seq:%d", id, seq)
		payload, err := (&icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: MagicCode,
			Body: &icmp.Echo{
				ID: id, Seq: seq,
				Data: wb[:size],
			},
		}).Marshal(nil)
		if err != nil {
			log.Fatalln("marshal echo error: ", err.Error())
		}
		_, err = serverConn.Write(payload)
		if err != nil {
			log.Fatalln("write echo request error: ", err.Error())
		}
		size = 0
		for {
			log.Println("wait client data")
			nr, _, err = serverConn.ReadFrom(rb)
			log.Println("get client data")
			if nr > 0 {
				wb = append(wb[:size], rb[:nr]...)
				size += nr
			}
			if err != nil {
				if err != io.EOF {
					log.Println("read server error: ", err)
					return
				}
			}

			reply, err := icmp.ParseMessage(ProtocolICMP, rb)
			if err != nil {
				log.Println("parse icmp echo reply error: ", err.Error())
				return
			}

			if reply.Code != MagicCode {
				return
			}

			body, err = reply.Body.Marshal(ProtocolICMP)
			if err != nil {
				log.Println("marshal icmp echo reply body error: ", err.Error())
				return
			}

			log.Printf("get echo reply id:%d and seq:%d",
				binary.BigEndian.Uint16(body[0:2]),
				binary.BigEndian.Uint16(body[2:4]))
			if binary.BigEndian.Uint16(body[0:2]) == uint16(id) &&
				binary.BigEndian.Uint16(body[2:4]) == uint16(seq) {
				log.Println("right")
				//ack
				ack, err := (&icmp.Message{
					Type: ipv4.ICMPTypeEcho, Code: MagicCode,
					Body: &icmp.Echo{
						ID: id, Seq: seq,
						Data: make([]byte, 0),
					},
				}).Marshal(nil)
				if err != nil {
					log.Println("marshal ack error:", err)
				}
				nw, err := serverConn.Write(ack)
				if err != nil {
					log.Println("write ack error", err)
				}
				log.Println("write ack size ", nw)
				break
			} else {
				log.Println("receive other")
				continue
			}
		}

		nr, err = clientConn.Write(body[4 : nr-4])
		if err != nil {
			log.Println("write client error: ", err.Error())
			return
		}
		log.Println("get echo reply size ", nr)
	}

}

func (self *HikarianIcmp) Run() {
	if self.mode == "decrypt" {
		clientConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			log.Fatal(err)
		}
		defer clientConn.Close()
		for {
			buf := make([]byte, 1024)
			nr, caddr, err := clientConn.ReadFrom(buf)
			request, err := icmp.ParseMessage(ProtocolICMP, buf)
			if err != nil {
				log.Println("parse icmp request error: ", err.Error())
				return
			}
			if request.Code == 0 {
				return
			}

			body, err := request.Body.Marshal(ProtocolICMP)
			if err != nil {
				log.Println("marshal body error: ", err.Error())
				continue
			}
			hash := binary.BigEndian.Uint16(body[0:2]) + binary.BigEndian.Uint16(body[2:4])
			channel := self.ChannelPool.Get(hash)
			if channel == nil {
				log.Println("new channel")
				channel = make(chan []byte)
				self.ChannelPool.Set(hash, channel)
				go self.transportServer(clientConn, caddr, channel)
			} else {
				log.Println("old channel")
			}
			channel <- body[:nr-4]
		}
	} else if self.mode == "encrypt" {
		client, err := net.ResolveTCPAddr("tcp", self.client)
		if err != nil {
			log.Fatalln("resolve client address error: ", err.Error())
		}
		l, err := net.ListenTCP("tcp", client)
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		for {
			clientConn, err := l.AcceptTCP()
			if err != nil {
				log.Println("accept error: ", err.Error())
				continue
			}
			defer clientConn.Close()
			go self.transportClient(clientConn)
		}
	}
}
