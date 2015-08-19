/*
* @Author: detailyang
* @Date:   2015-08-16 13:47:25
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-20 00:56:41
 */

package hikarian

import (
	"net"
)

type Conn struct {
	conn   net.Conn
	cipher *Cipher
}

func NewConn(conn net.Conn, cipher *Cipher) *Conn {
	return &Conn{
		conn:   conn,
		cipher: cipher,
	}
}

func (self *Conn) Read(b []byte) (int, error) {
	if self.cipher == nil {
		return self.conn.Read(b)
	}

	n, err := self.conn.Read(b)
	if n > 0 {
		self.cipher.Decrypt(b[0:n], b[0:n])
	}

	return n, err
}

func (self *Conn) Write(b []byte) (int, error) {
	if self.cipher == nil {
		return self.conn.Write(b)
	}

	self.cipher.Encrypt(b, b)
	return self.conn.Write(b)
}

func (self *Conn) Close() {
	self.conn.Close()
}

func (self *Conn) CloseRead() error {
	if conn, ok := self.conn.(*net.TCPConn); ok {
		return conn.CloseRead()
	}
	return nil
}

func (self *Conn) CloseWrite() error {
	if conn, ok := self.conn.(*net.TCPConn); ok {
		return conn.CloseRead()
	}
	return nil
}
