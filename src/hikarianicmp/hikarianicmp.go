package hikarianicmp

import (
	"encoding/binary"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
	"time"
)

const (
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
)

type HikarianIcmp struct {
	client, server string
	mode           string
}

func NewHikarianIcmp(client, server, mode string) *HikarianIcmp {
	return &HikarianIcmp{
		server: server,
		client: client,
		mode:   mode,
	}
}

func (self *HikarianIcmp) transportServer(clientConn *icmp.PacketConn) {
	for {
		buf := make([]byte, 1024)
		numRead, caddr, err := clientConn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			server, err := net.ResolveTCPAddr("tcp", self.server)
			if err != nil {
				log.Fatalln("resolve server address error: ", err.Error())
			}
			serverConn, err := net.DialTCP("tcp", nil, server)
			if err != nil {
				log.Println("connect remote address error:", err)
				return
			}
			serverConn.SetDeadline(time.Now().Add(time.Second * 5))
			defer serverConn.Close()

			request, err := icmp.ParseMessage(ProtocolICMP, buf)
			if err != nil {
				log.Println("parse icmp request error: ", err.Error())
				return
			}
			body, err := request.Body.Marshal(ProtocolICMP)
			if err != nil {
				log.Println("marshal body error: ", err.Error())
				return
			}

			_, err = serverConn.Write(body[4 : numRead-4])
			if err != nil {
				log.Println("write server error: ", err.Error())
				return
			}

			rb := make([]byte, 1024)
			wb := make([]byte, 1024)
			size := 0
			// for {
			nr, err := serverConn.Read(rb)
			if nr > 0 {
				size += nr
				wb = append(wb[size:], rb[:nr]...)
			}
			if err == io.EOF {
				// break
			}
			if err != nil {
				log.Println("read server error: ", err.Error())
				return
			}
			// }

			reply, err := (&icmp.Message{
				Type: ipv4.ICMPTypeEchoReply,
				Code: request.Code,
				Body: &icmp.Echo{
					ID:   int(binary.BigEndian.Uint16((body[0:2]))),
					Seq:  int(binary.BigEndian.Uint16((body[2:4]))),
					Data: wb[:size],
				},
			}).Marshal(nil)
			if err != nil {
				log.Println("marshal echo reply error: ", err.Error())
				return
			}
			numWrite, err := clientConn.WriteTo(reply, caddr)
			if err != nil {
				log.Println("write echo reply error: ", err.Error())
				return
			}
			numWrite = numWrite
		}()
	}
}

func (self *HikarianIcmp) transportClient(clientConn *net.TCPConn) {
	rb := make([]byte, 1024)
	wb := make([]byte, 1024)
	size := 0
	// for {
	nr, err := clientConn.Read(rb)
	if nr > 0 {
		log.Println(rb[:nr])
		size += nr
		wb = append(wb[:size], rb[:nr]...)
	}
	if err == io.EOF {
		// break
	}
	if err != nil {
		log.Println("read server error: ", err.Error())
		return
	}
	// }

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

	payload, err := (&icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: 0, Seq: 0,
			Data: wb,
		},
	}).Marshal(nil)
	if err != nil {
		log.Fatalln("marshal echo error: ", err.Error())
	}
	nw, err := serverConn.Write(payload)
	if err != nil {
		log.Fatalln("write echo request error: ", err.Error())
	}
	log.Println("here")
	size = 0
	// for {
	nr, err = serverConn.Read(rb)
	if nr > 0 {
		size += nr
		wb = append(wb[size:], rb[:nr]...)
	}
	if err == io.EOF {
		// break
	}
	if err != nil {
		log.Println("read server error: ", err.Error())
		return
	}
	// }

	nw, err = clientConn.Write(wb)
	if err != nil {
		log.Println("write client error: ", err.Error())
		return
	}

	log.Println("write size ", nw)
}

func (self *HikarianIcmp) Run() {
	if self.mode == "decrypt" {
		clientConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			log.Fatal(err)
		}
		defer clientConn.Close()
		self.transportServer(clientConn)
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
