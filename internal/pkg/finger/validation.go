package finger

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"sync"

	"Caesar/pkg/record"
)

// 单例模式变量
var (
	goInstance *AppFinger
	once       sync.Once
)

// Fingers 用来保存指纹数据
type Fingers []struct {
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
}

// AppFinger 指纹文件名
type AppFinger struct {
	fileName string
}

// 读取指纹
func (a AppFinger) getFingers() (Fingers, error) {
	var fingers Fingers

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(a.fileName)
	if err != nil {
		return nil, err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, &fingers)
	if err != nil {
		return nil, err
	}

	return fingers, nil

}

// CheckFinger 识别程序
func (a AppFinger) CheckFinger(html string) (result string, errs error) {
	fingersManager, err := a.getFingers()
	if err != nil {
		record.Logger.Error("Finger file read failed " + err.Error())
		return result, err
	}

	for _, v := range fingersManager {
		if strings.Contains(html, v.Keyword) {
			return v.Name, nil
		}

	}

	return result, nil

}

// NewLoads 单例模式生成对象
func NewLoads(fileName string) *AppFinger {

	if goInstance == nil {
		once.Do(func() {
			goInstance = &AppFinger{fileName: fileName}
		})
	}
	return goInstance

}
