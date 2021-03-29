package translate

import (
	"fmt"
	"github.com/bregydoc/gtranslate"
	"time"
)

func GoogleTranslate() {
	for {
		text := "西南科技大学坐落于中国科技城——四川省绵阳市。学校是四川省人民政府与教育部共建高校，四川省人民政府与国家国防科技工业局共建高校，被教育部确定为国家重点建设的西部14所高校之一。原中央政治局常委，国务院副总理李岚清同志赞誉学校“共建与区域产学研联合办学”走出了一条有自己特色的办学路子。学校现任党委书记陈永灿、校长董发勤。"
		translated, err := gtranslate.TranslateWithParams(
			text,
			gtranslate.TranslationParams{
				From: "zh",
				To:   "en",
			},
		)
		if err != nil {
			panic(err)
		}

		fmt.Printf("From: %s  \nTo: %s ", text, translated)
		time.Sleep(time.Second * 10)
	}

	// en: Hello World | ja: こんにちは世界
}
