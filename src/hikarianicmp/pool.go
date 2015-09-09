/*
* @Author: detailyang
* @Date:   2015-09-06 14:01:38
* @Last Modified by:   detailyang
* @Last Modified time: 2015-09-06 14:30:49
 */

package hikarianicmp

import "net"

type TCPConnPool struct {
	pool map[uint16]*net.TCPConn
}

func NewTCPConnPool() *TCPConnPool {
	return &TCPConnPool{
		pool: make(map[uint16]*net.TCPConn),
	}
}

func (self *TCPConnPool) Get(hash uint16) *net.TCPConn {
	return self.pool[hash]
}

func (self *TCPConnPool) Append(hash uint16, conn *net.TCPConn) {
	self.pool[hash] = conn
}
