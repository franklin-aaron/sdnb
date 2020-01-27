package main

import "scan/hostsurvival"

/**
 * @Description: 端口扫描工具
 * @File: main.go
 * @Time: 2020/1/19 16:55x
 */

//二级域名扫描
/*func checkDoman(url string) bool {

}*/
func main() {
	if err := hostsurvival.DomainScan("0430.com"); err != nil{
		panic(err)
	}
}
