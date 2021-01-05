package relation

import "time"

type EngineMap struct {
	// about console print message
	// value from config.yml
	Silence bool
	// Page similarity parameter
	// value from config.yml
	UpperRatioBound float64
	// request  threads
	// value from argv
	Threads int
	// check paths num
	//value from config.yml
	WafTop int
	// Record the number of timeout requests
	// value from config.yml
	TimeoutCount int
	// The delay between two requests
	// default: 0
	// value from argv
	TimeSleep int
	// The dir extend
	// such as zip, tar, rar
	// value from config.yml
	DirectoryDirSuffix []string
	// The dynamic(asp, php, jsp, aspx) extend
	// such as txt, bak, swp
	// value from config.yml
	DynamicFileSuffix []string
	// Display http code on result
	// such as 200, 302, 301
	// value from config.yml
	StatusCode []int
	// Will check path
	// value from assets/*.json
	//PathDict []TagPath
	// Stop scanning or not
	// value from dynamic running
	StopFlag bool
	// Symbolic link
	// such as: index.php.txt or index.php_txt
	// value from config.yml
	SuffixSymbol []string
	//Found Assets num
	// default: 0
	// value from dynamic running
	Numbers int
	// Scan result
	// value from dynamic running
	CollectAssets map[string][]ResultPtah
}

type PathsMap struct {
	//The program run dir
	BaseDir string
	// The path dir
	// default is : assets/directory
	DictDir string
	// The CDN  dir
	// default is : assets/cdn
	CdnPath string

	// The Finger  dir
	// default is : assets/fingerprint
	FingerPath string

	// Configuration file path
	// default is : configs/config.yml
	Config string
	// Result file save path
	// default is : results/result.json
	Result string

	// path map info
	// such as [spring:assets/directory/php.json]
	Dict []map[string]string
}

type BrowserMap struct {
	TimeOut   time.Duration
	UserAgent []string
	Proxy     string
	Cookie    string
	Headers   map[string]string
}
