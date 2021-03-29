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
	refer *string
}

var (
	UA        = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.3578.98 Safari/537.36"
	CookieStr = "tt_webid=6937485546063300126; ttcid=9f8a31d079e84b73988608c7786d5ff424; csrftoken=bcb1b35afeec7ebd4ed2c90992471963; tt_webid=6937485546063300126; __ac_nonce=060507c8a004d7fe7c338; __ac_signature=_02B4Z6wo00f01iOsgtQAAIDD9Gqvrp9OU0IjiIZAAOjJC1bjeblncBlXRoBQeRzsAB.dWew4JnFMWKLsrEJoJ4sl-p4bMowhvbNfQs.dzamie2mYhbvCCtcnaDKXGBkYDTAdAz398UdIBfLpb3; s_v_web_id=verify_kmbtp3du_tUj0yhVs_yflo_4iQZ_9NL9_9V6EzfW2WSDi; MONITOR_WEB_ID=fcef7a69-51f6-4ee0-b984-394f85b96268; tt_scid=4YkJx31sCEQ3KWkMQD2Nc2LiNMJEcPDjJyOFgNWVcMNRIiKSDNJ.mTwVmDsqcwvRc62a"
)

func NewHttpClient(proxyUrl ...string) *Req {
	req := new(Req)
	req.Client = http.Client{}
	req.Jar = new(Jar)
	req.Timeout = time.Second * 3
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	if len(proxyUrl) > 0 && proxyUrl[0] != "" {
		req.proxy = proxyUrl[0]
		proxyFunc := func(r *http.Request) (*url.URL, error) {
			r.Header.Set("User-Agent", UA)
			r.Header.Set("Referer", *req.refer)
			r.Header.Add("Cookie", CookieStr)
			return url.Parse(req.proxy)
		}
		tr.Proxy = proxyFunc
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	req.Transport = tr
	return req
}
func (req *Req) DoReq(method string, targetUrl string, callback func(*url.URL, io.Reader)) {
	switch method {
	case "POST", "GET", "PUT", "DELETE", "HEAD":
		urlInfo, err := url.Parse(targetUrl)
		if err != nil {
			log.Printf("can't parse url error %v", err)
		}
		if req.refer == nil {
			req.refer = &targetUrl
		}
		defer func() {
			req.refer = &targetUrl
		}()
		reqInfo, err := http.NewRequest(method, targetUrl, nil)
		if err != nil {
			log.Printf("new http request for url %s error", targetUrl)
			return
		}
		// no proxy
		if req.proxy == "" {
			reqInfo.Header.Set("User-Agent", UA)
			reqInfo.Header.Set("Referer", *req.refer)
			reqInfo.Header.Add("Cookie", CookieStr)
		}
		var resp *http.Response
		resp, err = req.Do(reqInfo)
		if err != nil {
			log.Printf("do request error %v", err)
			return
		}
		if resp.StatusCode == 200 {
			if callback != nil {
				callback(urlInfo, resp.Body)
			}
		} else {
			log.Printf("resp code %d not ok", resp.StatusCode)
		}
	default:
		log.Printf("no support method %s", method)
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
