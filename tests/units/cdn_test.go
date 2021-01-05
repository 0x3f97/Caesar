package units

import (
	"Caesar/internal/pkg/cdn"
	"Caesar/pkg/utils"
	"testing"
)

/*
检测CDN模块
*/
func TestCDN(t *testing.T) {
	f := cdn.NewIP("../../assets/cdn/cdn_ip_cidr.json")
	println(f.CheckIPCDN("192.168.3.1"))
	f2 := cdn.NewIP("../../assets/cdn/cdn_ip_cidr.json")
	println(f2.CheckIPCDN("192.168.3.1"))

}

func TestIP(t *testing.T) {
	println(utils.UrlToAddressAndPort("http://www.baidu.com:800"))
}
