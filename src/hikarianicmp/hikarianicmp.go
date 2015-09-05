package hikarianicmp

import (
	"encoding/binary"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
)

const (
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
)

type HikarianIcmp struct {
	client, server *net.TCPAddr
}

func NewHikarianIcmp(sserver string) *HikarianIcmp {
	server, err := net.ResolveTCPAddr("tcp", sserver)
	if err != nil {
		log.Fatal("resolve remote address failed")
	}

	return &HikarianIcmp{
		server: server,
	}
}

func (self *HikarianIcmp) transport(clientConn *icmp.PacketConn) {
	for {
		buf := make([]byte, 1024)
		numRead, caddr, err := clientConn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			serverConn, err := net.DialTCP("tcp", nil, self.server)
			if err != nil {
				log.Println("connect remote address error:", err)
				return
			}
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

			nw, err := serverConn.Write(body[4 : numRead-4])
			if err != nil {
				log.Println("write server error: ", err.Error())
				return
			}
			if nw != numRead {
				log.Println("write error")
			}

			rb := make([]byte, 1024)
			wb := make([]byte, 1024)
			size := 0
			for {
				nr, err := serverConn.Read(rb)
				if nr > 0 {
					size += nr
					wb = append(wb[size:], rb[:nr]...)
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Println("read server error: ", err.Error())
					break
				}
			}

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

func (self *HikarianIcmp) Run() {
	clientConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer clientConn.Close()
	self.transport(clientConn)
}
