/*
* @Author: detailyang
* @Date:   2015-08-23 17:51:43
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-23 17:52:20
 */

package hikariansocks5

const (
	HANDSHAKE  = 1 + 1 + 255
	VERSION    = 0x05
	NOAUTH     = 0x00
	NOACCEPT   = 0xFF
	TCPCONNECT = 0x01
	BIND       = 0x02
	UDP        = 0x03
	IPV4       = 0x01
	FQDN       = 0x03
	IPV6       = 0x04
	SUCCESS    = 0x00
	FAILURE    = 0x01
	BUFFSIZE   = 1024
)
