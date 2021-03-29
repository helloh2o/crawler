package validator

import (
	"CrawlerX/core"
	"io"
	"io/ioutil"
	"log"
	"net/url"
)

func Validate() {
	recaptchaGoogle := "https://www.google.com/recaptcha/api.js?render=6LfRG68UAAAAAMToh2v5n7aqfEyrhVD584F8JL20&amp;onload=captchaReady&amp;hl=zh-tw"
	req := core.NewHttpClient("http://103.135.250.118:38380")
	req.DoReq("GET", recaptchaGoogle, func(url *url.URL, reader io.Reader) {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		log.Println(string(data))
	})
}
