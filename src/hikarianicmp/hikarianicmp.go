package hikarianicmp

import (
	"encoding/binary"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
)

const (
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
)

type HikarianIcmp struct {
}

func NewHikarianIcmp() *HikarianIcmp {
	return &HikarianIcmp{}
}

func (self *HikarianIcmp) transport(clientConn *icmp.PacketConn) {
	buf := make([]byte, 1024)
	for {
		numRead, caddr, err := clientConn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		request, err := icmp.ParseMessage(ProtocolICMP, buf)
		if err != nil {
			log.Println("parse icmp request error: ", err.Error())
			continue
		}
		body, err := request.Body.Marshal(ProtocolICMP)
		if err != nil {
			log.Println("marshal body error: ", err.Error())
		}
		reply, err := (&icmp.Message{
			Type: ipv4.ICMPTypeEchoReply,
			Code: request.Code,
			Body: &icmp.Echo{
				ID:   int(binary.BigEndian.Uint16((body[0:2]))),
				Seq:  int(binary.BigEndian.Uint16((body[2:4]))),
				Data: body[4 : numRead-4],
			},
		}).Marshal(nil)
		if err != nil {
			log.Println("marshal echo reply error: ", err.Error())
		}
		numWrite, err := clientConn.WriteTo(reply, caddr)
		if err != nil {
			log.Println("write echo reply error: ", err.Error())
		}
		numWrite = numWrite
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
