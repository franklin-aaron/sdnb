package main

import "scan/hostsurvival"

/**
 * @Description: 端口扫描工具
 * @File: main.go
 * @Time: 2020/1/19 16:55x
 */

func main() {
	if err := hostsurvival.DomainScan("google.com"); err != nil{
		panic(err)
	}
}
