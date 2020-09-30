package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"regexp"
	"os"
	"flag"
	json "github.com/tidwall/gjson"
	"fmt"
	"strconv"

)

func usage() {
    print("\n" + appname + "/" + version + "\n用法：" + os.Args[0] + " -l Length SongName\n参数：\n")
	flag.PrintDefaults()
	print("\n")
}

var appname string = "migu.get"
var version string = "1.0.0g"

func main()  {
	var h, v bool
	var l int
	flag.IntVar(&l, "l", 12, "`Length`：搜索结果的长度")
	flag.BoolVar(&v, "v", false, "显示版本")
	flag.BoolVar(&h, "h", false, "显示帮助")
	flag.Parse()
	if h || len(os.Args) == 1 {
		flag.Usage = usage
		flag.Usage()
		os.Exit(0)
	}
	if v {
		print(appname + "/" + version + "\n")
		os.Exit(0)
	}
	var R1 Rr
	sS := `{"song":1,"album":0,"singer":0,"tagSong":1,"mvSong":0,"songlist":0,"bestShow":1,"lyricSong":0,"concert":0,"periodical":0,"ticket":0,"bit24":0,"verticalVideoTone":0}`
	R1.u = "http://jadeite.migu.cn:7090/music_search/v2/search/searchAll?isCopyright=1&isCorrect=1&pageIndex=1&pageSize=" + strconv.Itoa(l) + "&searchSwitch=" + url.QueryEscape(sS) + "&text=" + url.QueryEscape(os.Args[len(os.Args)-1])
	R1.h = [][]string{
		{"timeStamp", "1596391260"},
		{"sign", "98d131a54422139907d45f7f204ecf72"},
		{"version", "666.666.666"},
	}
	r1 := Http(R1)
	R := json.Get(r1.b, "songResultData.result").Array()
	if len(R) == 0 {
		print("搜不到...\n")
		os.Exit(0)
	}
	print("“" + os.Args[len(os.Args)-1] + "”的搜索结果：\n============\n")
	for i := 0; i < len(R); i++ {
		Sn := R[i].Get("songName").String()
		Sg := R[i].Get("singer").String()
		Al := R[i].Get("album").String()
		n := strconv.Itoa(i+1)
		if i < 9 {
			n = strconv.Itoa(i+1) + " "
		}
		print(n + "  " + Sn + " - " + Sg + " - " + Al + "\n")
	}
	print("============\n请输入歌曲序号：")
	r := ""
	fmt.Scan(&r)
	w, err := strconv.Atoi(r)
	if err != nil {
		print("输入有误!!\n")
		os.Exit(0)
	}
	if len(R) < w-1 {
		print("超出范围!!\n")
		os.Exit(0)
	}
	Su := R[w-1].Get("newRateFormats").Array()
	if len(Su) == 0 {
		print("nothing...\n")
		os.Exit(0)
	}
	for i := 0; i < len(Su); i++ {
		Ft := Su[i].Get("formatType").String()
		if Ft == "SQ" || Ft == "ZQ" {
			u := match(Su[i].Get("androidUrl").String(), `ftp://.*?/`, "http://freetyst.nf.migu.cn/", "")
			print(Ft + ":\n" + u[0][0] + "\n")
		}else{
			u := match(Su[i].Get("url").String(), `ftp://.*?/`, "http://freetyst.nf.migu.cn/", "")
			print(Ft + ":\n" + u[0][0] + "\n")
		}
	}
}

type Rr struct {
	Url,u string
	Header,h [][]string
	Status,s string
	Data,d string
	Body,b string
	Method,m string
	ReturnHeader,rh string
}

func Http(Req Rr) Rr{
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
			},
	}
	D := strings.NewReader(Req.d)
	request, err := http.NewRequest(Req.m, Req.u, D)
	if(len(Req.h) != 0){
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for i := 0; i < len(Req.h); i++ {
			request.Header.Set(Req.h[i][0], Req.h[i][1])
		}
	}
	if err != nil {
		print("ERROR!!0x22\n")
		os.Exit(0)
	}
	response, err := client.Do(request)
	if err != nil {
		print("ERROR!!0x23\n")
		os.Exit(0)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		print("ERROR!!0x33\n")
		os.Exit(0)
	}
	h := response.Header
	hs := ""
	for k := range h {
		hs = hs + k + ":" + h[k][0] + "\r\n"
	}
	var ret Rr
	ret.rh = hs + response.Status + "\r\n"
	ret.b = string(body)
	ret.s = response.Status
	return ret
}

func match(nr, reg, rep, err string) [][]string{
	p := regexp.MustCompile(reg)
	if rep == "" {
		result := p.FindAllStringSubmatch(nr, -1)
		if len(result) != 0 {
			return result
		}else{
			if err != "" {
				print(err + "\n")
				os.Exit(0)
			}
			return nil
		}		
	}else{
		result := p.ReplaceAllString(nr, rep)
		return [][]string{{result}}
	}
}
