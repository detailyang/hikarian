/*
* @Author: detailyang
* @Date:   2015-08-21 21:44:50
* @Last Modified by:   detailyang
* @Last Modified time: 2015-08-23 19:11:17
 */

package main

import (
	"hikariansocks5"
)

func main() {
	hs := hikariansocks5.NewHikarianSocs5(":1080")

	hs.Run()
}
