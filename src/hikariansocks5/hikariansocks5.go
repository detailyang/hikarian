/*
* @Author: detailyang
* @Date:   2015-08-21 00:12:01
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-23 21:52:40
 */

package hikariansocks5

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

type HikarianSocks5 struct {
	server *net.TCPAddr
}

func NewHikarianSocs5(server string) *HikarianSocks5 {
	s, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		log.Fatalln("resolve server address failed ", err)
	}

	return &HikarianSocks5{
		server: s,
	}
}

func (self *HikarianSocks5) handshake(s *Socks5) error {
	_, err := s.cconn.Read(s.handshake)
	if err != nil {
		return err
	}

	/*
	   +----+----------+----------+
	   |VER | NMETHODS | METHODS  |
	   +----+----------+----------+
	   | 1  |    1     | 1 to 255 |
	   +----+----------+----------+
	*/
	/*
	   0x00        NO AUTHENTICATION REQUIRED(无需认证)
	   0x01        GSSAPI
	   0x02        USERNAME/PASSWORD(用户名/口令认证机制)
	   0x03-0x7F   IANA ASSIGNED
	   0x80-0xFE   RESERVED FOR PRIVATE METHODS(私有认证机制)
	   0xFF        NO ACCEPTABLE METHODS(完全不兼容)
	*/
	if s.handshake[0] != VERSION {
		return errors.New("only support socks5")
	}

	if s.handshake[2] != NOAUTH {
		return errors.New("only support no auth")
	}

	return nil
}

func (self *HikarianSocks5) rejectHandshake(s *Socks5) (int, error) {
	return s.cconn.Write([]byte{VERSION, NOACCEPT})
}

func (self *HikarianSocks5) acceptHandshake(s *Socks5) (int, error) {
	return s.cconn.Write([]byte{VERSION, NOAUTH})
}

func (self *HikarianSocks5) requests(s *Socks5) error {
	/*
	   +----+-----+-------+------+----------+----------+
	   |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	   +----+-----+-------+------+----------+----------+
	   | 1  |  1  | X'00' |  1   | Variable |    2     |
	   +----+-----+-------+------+----------+----------+
	*/
	n, err := s.cconn.Read(s.requests)
	s.requests_size = n
	if err != nil {
		return err
	}

	if s.requests[1] != TCPCONNECT {
		return errors.New("only support tcp connect")
	}

	return nil
}

func (self *HikarianSocks5) rejectRequests(s *Socks5) (int, error) {
	/*
	   +----+-----+-------+------+----------+----------+
	   |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	   +----+-----+-------+------+----------+----------+
	   | 1  |  1  | X'00' |  1   | Variable |    2     |
	   +----+-----+-------+------+----------+----------+
	*/
	s.requests[1] = FAILURE
	return s.sconn.Write(s.requests[:s.requests_size])
}

func (self *HikarianSocks5) replies(s *Socks5) error {
	ip := "0.0.0.0"
	port := uint16(80)
	size := 0

	switch s.requests[3] {
	case IPV4:
		log.Println("receive IPV4 request")
		ip = net.IPv4(s.requests[4], s.requests[5], s.requests[6], s.requests[7]).String()
		port = binary.BigEndian.Uint16(s.requests[8:10])
		s.replies[3] = IPV4
		size = 10
	case IPV6:
		log.Println("receive IPV6 request")
		s.replies[3] = IPV6
	case FQDN:
		log.Println("receive FQDN request")
		ip = string(s.requests[5 : 5+s.requests[4]])
		port = binary.BigEndian.Uint16(s.requests[5+s.requests[4]:])
		s.replies[3] = FQDN
		copy(s.replies[4:], s.requests[4:])
		size = 5 + int(s.requests[4]) + 2
	default:
		return errors.New("unknow requests type")
	}
	s.replies[0] = VERSION
	s.replies[1] = SUCCESS
	s.replies[2] = 0x00

	sconn, err := net.Dial("tcp", ip+":"+strconv.Itoa(int(port)))
	s.sconn = sconn
	log.Println("connect address success: ", ip+":"+strconv.Itoa(int(port)))
	if err != nil {
		s.replies[1] = FAILURE
		s.cconn.Write(s.replies[:size])
		return err
	}
	_, err = s.cconn.Write(s.replies[:size])
	if err != nil {
		return err
	}

	return nil
}

func (self *HikarianSocks5) pipe(dst net.Conn, src net.Conn) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Println("io copy error ", err)
		dst.Close()
		src.Close()
	}
}

func (self *HikarianSocks5) transport(conn net.Conn) {
	socks5 := NewSocks5(conn)
	//handshake
	err := self.handshake(socks5)
	if err != nil {
		log.Println("hankshake failed: ", err)
		_, err := self.rejectHandshake(socks5)
		if err != nil {
			log.Println("reject hankshake failed: ", err)
		}
		return
	}
	_, err = self.acceptHandshake(socks5)
	if err != nil {
		log.Println("accept handshake failed: ", err)
		return
	}
	log.Println("accept handshake success: ", socks5.cconn.RemoteAddr())

	//requests
	err = self.requests(socks5)
	if err != nil {
		log.Println("requests failed: ", err)
		_, err := self.rejectRequests(socks5)
		if err != nil {
			log.Println("reject requests failed: ", err)
		}
		return
	}

	err = self.replies(socks5)
	if err != nil {
		log.Println("replies failed: ", err)
		return
	}

	go self.pipe(socks5.sconn, socks5.cconn)
	go self.pipe(socks5.cconn, socks5.sconn)
}

func (self *HikarianSocks5) Run() {
	l, err := net.ListenTCP("tcp", self.server)
	if err != nil {
		log.Fatalln("listen server address failed ", err)
	}
	defer l.Close()
	log.Println("listen on ", self.server)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept failed ", err)
			continue
		}
		go self.transport(conn)
	}
}
