/*
* @Author: detailyang
* @Date:   2015-09-06 14:01:38
* @Last Modified by:   detailyang
* @Last Modified time: 2015-09-13 00:44:29
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

func (self *TCPConnPool) Set(hash uint16, conn *net.TCPConn) {
	self.pool[hash] = conn
}

type ChannelPool struct {
	pool map[uint16]chan []byte
}

func NewChannelPool() *ChannelPool {
	return &ChannelPool{
		pool: make(map[uint16]chan []byte),
	}
}

func (self *ChannelPool) Get(hash uint16) chan []byte {
	return self.pool[hash]
}

func (self *ChannelPool) Set(hash uint16, channel chan []byte) {
	self.pool[hash] = channel
}
