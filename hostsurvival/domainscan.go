package hostsurvival

/**
 * @Description: 二级域名爆破
 * @File: domainscan.go
 * @Time: 2020/1/22 18:29
 */
import (
	"fmt"
	"github.com/modood/table"
	"net"
	"runtime"
	"scan/file"
)

type DomainDns struct {
	Ip     [] net.IP
	Cname  string
	Ptr    [] string
	Ns     [] *net.NS
	Mx     [] *net.MX
	Srv    Srv
	Txt    [] string
	domain string //域名
}

type Srv struct {
	cname string
	addrs [] *net.SRV
}

type House struct {
	Ip     string
	Cname  string
	Ptr    string
	Ns     string
	Mx     string
	Txt    string
	Domain string //域名
}

//进行检查
func checkDomain(domain string, domainNamePrefix chan string, exitChan chan bool, data chan DomainDns) {

	for {
		v, ok := <-domainNamePrefix
		if !ok {
			break
		}
		if v != domain {
			v = fmt.Sprintf("%v.%v", v, domain)
		}

		ips, _ := net.LookupIP(v)

		//判定该主机是否有解析
		if len(ips) != 0 {
			//域名的cname
			cname, err := net.LookupCNAME(v)
			if err != nil {
				panic(err)
			}

			//域名的ptr
			ptr := make([] string, 0, 50)
			for _, ip := range ips {
				p, err := net.LookupAddr(ip.String())
				if err != nil {
					break
				}
				ptr = append(ptr, p...)
			}

			//域名ns
			ns, _ := net.LookupNS(v)

			//域名MX
			mx, _ := net.LookupMX(v)

			//域名srv
			srvCname, addrs, _ := net.LookupSRV("xmpp-server", "tcp", v)
			srvVal := Srv{
				cname: srvCname,
				addrs: addrs,
			}

			//域名txt记录
			txt, _ := net.LookupTXT(v)

			rsf := DomainDns{
				Ip:     ips,
				Cname:  cname,
				Ptr:    ptr,
				Ns:     ns,
				Mx:     mx,
				Srv:    srvVal,
				Txt:    txt,
				domain: v,
			}

			data <- rsf
		}
	}

	//进程完毕,写入关闭管道
	exitChan <- true
}

//使用字典对域名进行扫描
func DomainScan(domain string) error {
	//存放结果
	data := make(chan DomainDns, 50)

	//打开字典
	reader, err := file.GetInstance("./dic/domain.dic")
	if err != nil {
		fmt.Println(err)
		return err
	}

	//逐行读取， 支持以切片或管道方式
	domainNamePrefix := make(chan string, 1500)
	domainNamePrefix <- domain //将顶级域名放入
	if err := reader.ReadLine(nil, domainNamePrefix); err != nil {
		fmt.Println(err)
		return err
	}

	//协程检查域名
	cpuNum := runtime.NumCPU()          //设定协程数量
	exitChan := make(chan bool, cpuNum) //关闭管道
	for i := 0; i <= cpuNum; i++ {
		go checkDomain(domain, domainNamePrefix, exitChan, data)
	}

	//关闭结果管道
	go func() {
		for i := 0; i <= cpuNum; i++ {
			_,ok := <-exitChan
			if !ok {
				break
			}
		}
		close(data)
	}()

	s := make([] House, 0, 50)

	//读取结果
	for {
		rsf, ok := <-data
		if !ok {
			break
		}

		//将ip 转为字符串
		ips := ""
		for _, ipV := range rsf.Ip {
			ips += ipV.String() + "|"
		}

		//将ptr转为字符
		ptrs := ""
		for _,ptrV := range rsf.Ptr{
			ptrs += ptrV
		}

		//将ns转为字符串
		nss := ""
		for _, ns := range rsf.Ns {
			nss += (*ns).Host + "|"
		}

		//将mx转为字符
		mxs := ""
		for _, mx := range rsf.Mx {
			mxs += mx.Host + " 优先级：" + string(mx.Pref) + "|"
		}

		//将txt转化
		txtrecords := ""
		for _, txt := range rsf.Txt {
			txtrecords += txt + "|"
		}

		s = append(s, House{
			Ip:     ips,
			Cname:  rsf.Cname,
			Ptr:    ptrs,
			Ns:     nss,
			Mx:     mxs,
			Txt:    txtrecords,
			Domain: rsf.domain,
		}, House{})
	}

	//表格终端输出
	t := table.Table(s)
	fmt.Println(t)
	
	return nil
}
