package core

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Req struct {
	http.Client
	proxy string
}

var (
	agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.3578.98 Safari/537.36"
)

func NewReq(proxyUrl string) *Req {
	req := new(Req)
	req.Client = http.Client{}
	req.Jar = new(Jar)
	req.Timeout = time.Second * 3
	tr := &http.Transport{}
	if proxyUrl != "" {
		proxyFunc := func(r *http.Request) (*url.URL, error) {
			r.Header.Set("User-Agent", agent)
			return url.Parse(proxyUrl)
		}
		tr.Proxy = proxyFunc
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req.Transport = tr
	return req
}

func (req *Req) Crawl(targetUrl string, callback func(*url.URL, io.Reader)) {
	urlInfo, err := url.Parse(targetUrl)
	if err != nil {
		log.Printf("can't parse url error %v", err)
	}
	reqInfo, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Printf("new http request for url %s error", targetUrl)
		return
	}
	reqInfo.Header.Add("Host", urlInfo.Host)
	reqInfo.Header.Add("User-Agent", agent)
	resp, err := req.Do(reqInfo)
	if err != nil {
		log.Printf("do request error %v", err)
		return
	}
	if resp.StatusCode == 200 {
		if callback != nil {
			callback(urlInfo, resp.Body)
		}
	}
}

type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}