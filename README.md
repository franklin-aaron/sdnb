# sdnb
二级域名爆破，基于go 语言实现，字典是在github 上所获得

go version : > 1.8

usage of application:

下载代码 ： git clone https://github.com/franklin-aaron/sdnb.git

进入main.go ： 
 例如：（不要加前缀） 
    func main() {
      if err := hostsurvival.DomainScan("google.com"); err != nil{
        panic(err)
      }
    }

效果如下： 
![image](https://github.com/franklin-aaron/sdnb/blob/master/photo_2020-02-02_17-02-58.jpg)
