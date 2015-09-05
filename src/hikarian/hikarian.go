/*
* @Author: detailyang
* @Date:   2015-08-16 11:03:47
* @Last Modified by:   detailyang
* @Last Modified time: 2015-09-05 14:49:39
 */

package hikarian

import (
	"io"
	"log"
	"net"
)

type Hikarian struct {
	client, server *net.TCPAddr
	mode           string
	algo           string
	secret         string
}

func NewHikarian(client, server, mode, algo, secret string) *Hikarian {
	c, err := net.ResolveTCPAddr("tcp", client)
	if err != nil {
		log.Fatalln("resolve client address error:", err)
	}

	s, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		log.Fatalln("resolve client address error:", err)
	}

	return &Hikarian{
		client: c,
		server: s,
		mode:   mode,
		secret: secret,
		algo:   algo,
	}
}

func (self *Hikarian) transport(clientConn *net.TCPConn) {
	serverConn, err := net.DialTCP("tcp", nil, self.server)
	if err != nil {
		log.Println("connect remote address error:", err)
		return
	}
	cipher, err := NewChiper(self.algo, self.secret)
	if err != nil {
		log.Fatalf("generate chiper: %s failed %s", self.algo, err.Error())
		return
	}

	var c, s *Conn
	if self.mode == "encrypt" {
		c = NewConn(clientConn, nil)
		s = NewConn(serverConn, cipher)
	} else if self.mode == "decrypt" {
		c = NewConn(clientConn, cipher)
		s = NewConn(serverConn, nil)
	} else {
		c = NewConn(clientConn, nil)
		s = NewConn(serverConn, nil)
	}

	go self.pipe(c, s)
	go self.pipe(s, c)
}

func (self *Hikarian) pipe(src, dst *Conn) {
	defer func() {
		src.CloseRead()
		dst.CloseWrite()
	}()
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Println("copy src to dst error:", err)
	}
}

func (self *Hikarian) Run() {
	l, err := net.ListenTCP("tcp", self.client)
	if err != nil {
		log.Fatalln("listen error:", err)
	}
	defer l.Close()
	log.Println("listen to client ", self.client)

	for {
		clientConn, err := l.AcceptTCP()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go self.transport(clientConn)
	}
}
