/*
* @Author: detailyang
* @Date:   2015-09-05 15:17:25
* @Last Modified by:   detailyang
* @Last Modified time: 2015-09-05 17:38:18
 */

/*
package main

import "log"
import "golang.org/x/net/icmp"
import "golang.org/x/net/ipv4"
import "net"
import "os"

func main() {
    laddr := &net.IPAddr{IP: net.ParseIP("0.0.0.0")}
    raddr, err := net.ResolveIPAddr("ip", os.Args[1])
    if err != nil {
        log.Fatalln("parse remote addr error: ", err.Error())
    }

    conn, err := net.DialIP("ip4:icmp", laddr, raddr)
    if err != nil {
        log.Fatalln("dial ip failed", err.Error())
        return
    }
    defer conn.Close()

    payload, err := (&icmp.Message{
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID: 0, Seq: 0,
            Data: []byte("abcd"),
        },
    }).Marshal(nil)
    if err != nil {
        log.Fatalln("marshal echo error: ", err.Error())
    }
    log.Println(payload)
    numWrite, err := conn.Write(payload)
    if err != nil {
        log.Fatalln("write echo request error: ", err.Error())
    }
    log.Println("write size: ", numWrite)
    for {
    }
}

*/