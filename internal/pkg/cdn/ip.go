package cdn

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"sync"

	"Caesar/pkg/record"
)

var (
	goInstance *IPFile
	once       sync.Once
)

type ipCDN []string

// IPFile 是读取CDN ip段的信息
type IPFile struct {
	Name string
}

func (ifc *IPFile) getIP() (ipCDN, error) {
	// 读取json文件
	var info ipCDN

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(ifc.Name)
	if err != nil {
		return nil, err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// CheckIPCDN 用来检查IP地址是否是CDN
func (ifc *IPFile) CheckIPCDN(ip string) bool {

	ipRange, err := ifc.getIP()
	if err != nil {
		record.Logger.Error("CDN ip file read failed " + err.Error())
		return false
	}

	ipFormat := net.ParseIP(ip)
	for _, v := range ipRange {
		_, ipNet, _ := net.ParseCIDR(v)
		if ipNet.Contains(ipFormat) {
			return true
		}

	}

	return false
}

// NewIP 使用go 实现单例模式
func NewIP(name string) *IPFile {
	if goInstance == nil {
		once.Do(func() {
			goInstance = &IPFile{
				Name: name,
			}
		})
	}
	return goInstance
}
