package boot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"Caesar/internal/library"
	"Caesar/internal/relation"
	"Caesar/pkg/utils"
)

func SetPaths() {
	base, _ := os.Getwd()
	var dictMap []map[string]string

	relation.Paths.BaseDir = base

	// 设置字典路径
	relation.Paths.DictDir = filepath.Join(base, "assets", "directory")

	// 设置CDN数据路径
	relation.Paths.CdnPath = filepath.Join(base, "assets", "cdn")

	// 设置web指纹数据路径
	relation.Paths.FingerPath = filepath.Join(base, "assets", "fingerprint")

	// 设置配置文件路径
	if !utils.PathExists(filepath.Join(base, "configs", "config.yml")) {
		if !utils.PathExists(filepath.Join(base, "config.yml")) {
			println("Can not found config file: " + filepath.Join(base, "config.yml"))
			os.Exit(1)

		} else {
			relation.Paths.Config = filepath.Join(base, "config.yml")
		}

	} else {
		relation.Paths.Config = filepath.Join(base, "configs", "config.yml")
	}

	// 设置结果保存路径
	relation.Paths.Result = filepath.Join(base, "results")

	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(relation.Paths.DictDir)

	if err != nil {
		println("The dir \""+relation.Paths.DictDir+"\" read error ", err.Error())
		os.Exit(1)
	}

	for _, v := range fileInfoList {
		key := strings.Split(v.Name(), ".json")[0]
		dictMap = append(dictMap, map[string]string{key: relation.Paths.DictDir + "/" + v.Name()})
	}

	relation.Paths.Dict = dictMap

}

func SetConf() {

	config := library.NewProfile(relation.Paths.Config).LoadConfigFromYaml()

	// 核心引擎设置
	relation.Engine.TimeSleep = 0
	relation.Engine.Silence = config.Silence
	relation.Engine.UpperRatioBound = config.UpperRatioBound

	// 设置默认线程数
	relation.Engine.Threads = 3

	relation.Engine.WafTop = config.WafTop
	relation.Engine.TimeoutCount = config.TimeoutCount
	relation.Engine.DynamicFileSuffix = config.DirectorySuffix
	relation.Engine.DirectoryDirSuffix = config.DynamicSuffix
	relation.Engine.SuffixSymbol = config.SuffixConnector
	relation.Engine.StatusCode = config.DisplayCode
	relation.Engine.StopFlag = false
	relation.Engine.CollectAssets = make(map[string][]relation.ResultPtah)

	// 浏览器参数设置
	relation.Browser.TimeOut = time.Duration(config.TimeOut)

	if len(config.UserAgent) > 0 {
		relation.Browser.UserAgent = config.UserAgent
	} else {
		relation.Browser.UserAgent = []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:76.0) Gecko/20100101 Firefox/76.0"}
	}

	if len(config.Proxy) > 0 {
		relation.Browser.Proxy = config.Proxy
	}

	if len(config.Cookie) > 0 {
		relation.Browser.Cookie = config.Cookie
	} else {
		relation.Browser.Cookie = ""
	}

	if len(config.Headers) > 0 {
		relation.Browser.Headers = config.Headers
	} else {
		relation.Browser.Headers = make(map[string]string)
	}

}
