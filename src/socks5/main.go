/*
* @Author: detailyang
* @Date:   2015-08-21 21:44:50
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-23 21:22:00
 */

package main

import (
	"flag"
	"hikariansocks5"
)

func main() {
	var server string
	flag.StringVar(&server, "server", "127.0.0.1:1080", "server listen")
	flag.Parse()
	hs := hikariansocks5.NewHikarianSocs5(server)

	hs.Run()
}
