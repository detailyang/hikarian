/*
* @Author: detailyang
* @Date:   2015-08-23 17:31:33
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-23 18:22:33
 */

package hikariansocks5

import (
	"net"
)

type Socks5 struct {
	cconn, sconn                 net.Conn
	handshake, requests, replies []byte
	requests_size                int
}

func NewSocks5(cconn net.Conn) *Socks5 {
	return &Socks5{
		cconn:         cconn,
		sconn:         nil,
		handshake:     make([]byte, HANDSHAKE),
		requests:      make([]byte, BUFFSIZE),
		replies:       make([]byte, BUFFSIZE),
		requests_size: 0,
	}
}
