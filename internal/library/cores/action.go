package cores

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"Caesar/internal/library/director"
	"Caesar/internal/library/engine"
	"Caesar/internal/library/extra"
	"Caesar/internal/pkg/cdn"
	"Caesar/internal/relation"
	"Caesar/internal/report"
	"Caesar/pkg/builder/generated"
	"Caesar/pkg/record"
	"Caesar/pkg/utils"
)

func Start(target string, typeRequest bool, paths []relation.TagPath) {
	/*
		typeRequest:true -> target是url地址
					false -> target是一个标准的http请求文本
	*/

	var enablePaths []relation.TagPath
	var threads int
	var req engine.RequestInfo

	var waf = false
	var mvc = false

	if typeRequest {
		req = engine.RequestInfo{
			Address:   target,
			Method:    http.MethodGet,
			Header:    relation.Browser.Headers,
			Cookies:   relation.Browser.Cookie,
			Body:      "",
			Proxy:     relation.Browser.Proxy,
			UserAgent: relation.Browser.UserAgent,
			Timeout:   relation.Browser.TimeOut,
		}

		// 检查目标连通性
		if _, _, _, err := director.UnitTest(req.Address, req.Method, utils.GetRandomElement(req.UserAgent), req.Header, req.Proxy, req.Cookies, relation.Browser.TimeOut, req.Body); err != nil {
			record.Logger.Fatal("Can not connect " + target)
			return
		}
	} else {

		target, method, agent, cookies, headers, data := generated.ParseRequestFromFile(target)

		req = engine.RequestInfo{
			Address:   utils.ConvertAddress(target),
			Method:    method,
			Header:    headers,
			Cookies:   cookies,
			Body:      data,
			Proxy:     relation.Browser.Proxy,
			UserAgent: []string{agent},
			Timeout:   relation.Browser.TimeOut,
		}

		// 检查目标连通性
		if _, _, _, err := director.UnitTest(req.Address, req.Method, utils.GetRandomElement(req.UserAgent), req.Header, req.Proxy, req.Cookies, relation.Browser.TimeOut, req.Body); err != nil {
			record.Logger.Fatal("Can not connect " + req.Address)
			return
		}
	}

	// 检查不存在页面返回的信息
	if status, header, body, err := director.UnitTest(req.Address+"/"+utils.GenRandString(6), req.Method, utils.GetRandomElement(req.UserAgent), req.Header, req.Proxy, req.Cookies, relation.Browser.TimeOut, req.Body); err == nil {

		// 提取目标地址
		address, _ := utils.UrlToAddressAndPort(req.Address)
		ip := utils.DomainToIP(address)

		// 测试是否是内网
		if utils.IsPrivateIP(ip) {
			record.Logger.Warn("The target is Private " + req.Address)

		} else {
			// 检查是否是cdn
			if cdn.NewIP(relation.Paths.CdnPath + "/cdn_ip_cidr.json").CheckIPCDN(ip) {
				record.Logger.Warn("The target is CDN " + req.Address)
			}
		}

		// 检查是否存在WAF
		if CheckWaf(req.Address) {
			waf = true
			record.Logger.Warn("The target by waf project " + req.Address)
			// 睡3秒
			time.Sleep(3 * time.Second)

		}

		// 检查是否是MVC框架
		if checkFrame, frame := CheckMVC(body); checkFrame {
			mvc = true
			record.Logger.Warn("The target is " + frame)

		}

		threads = relation.Engine.Threads

		enablePaths = paths

		if len(paths) > relation.Engine.WafTop {
			enablePaths = extra.GetFilterPath(paths, mvc, relation.Engine.WafTop)
		}

		// 开始构造参数
		resp := engine.ResponseInfo{
			Header: header,
			Body:   body,
		}

		opts := engine.ServerOpt{
			Paths:   enablePaths,
			Threads: threads,
			WAF:     waf,
		}

		if mvc {
			gun := engine.CreateFactory(status, req, resp, opts)
			engine.MVCFuzz(gun)
		} else {
			gun := engine.CreateFactory(status, req, resp, opts)
			engine.StandFuzz(gun)
		}

		return

	}

}

func Export(results map[string][]relation.ResultPtah) {
	// 最终结果导出函数
	var resultsList []relation.ResultPtah

	for key, value := range results {

		record.Logger.Debug(key + " found " + strconv.Itoa(len(value)) + " assets")

		if !relation.Engine.Silence {
			println()
			// 善后工作，在终端打印结果
			for _, v := range value {
				println(v.Address + " " + v.Title + " " + strconv.Itoa(v.Code) + " " + strconv.Itoa(v.Length))
			}
			println()
		}

		resultsList = append(resultsList, value...)

	}

	if len(resultsList) > 0 {
		report.ExportJson(resultsList, filepath.Join(relation.Paths.Result, "results.json"))

	}
}
