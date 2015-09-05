/*
* @Author: detailyang
* @Date:   2015-09-05 17:00:56
* @Last Modified by:   detailyang
* @Last Modified time: 2015-09-05 20:42:23
 */

package main

import (
	"flag"
	"hikarianicmp"
)

func main() {
	var client, server, mode, secret, algo string
	flag.StringVar(&client, "client", "127.0.0.1:3000", "client listen")
	flag.StringVar(&server, "server", "127.0.0.1:6378", "connect to server")
	flag.StringVar(&mode, "mode", "none", "mode")
	flag.StringVar(&algo, "algo", "rc4", "algo")
	flag.StringVar(&secret, "secret", "iamyoufather", "encrypt secret")
	flag.Parse()

	hi := hikarianicmp.NewHikarianIcmp(server)
	hi.Run()
}
