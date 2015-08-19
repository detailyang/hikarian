/*
* @Author: detailyang
* @Date:   2015-08-12 01:10:50
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-20 01:05:35
 */

package main

import (
	"flag"
	"hikarian"
)

func main() {
	var client, server, mode, secret, algo string
	flag.StringVar(&client, "client", "127.0.0.1:3000", "client listen")
	flag.StringVar(&server, "server", "127.0.0.1:6378", "connect to server")
	flag.StringVar(&mode, "mode", "none", "mode")
	flag.StringVar(&algo, "algo", "rc4", "algo")
	flag.StringVar(&secret, "secret", "iamyoufather", "encrypt secret")
	flag.Parse()

	h := hikarian.NewHikarian(client, server, mode, algo, secret)
	h.Run()
}
