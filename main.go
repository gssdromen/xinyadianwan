package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"strconv"

	"github.com/0xAX/notificator"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

func main() {
	// var HomePageUrl = "http://www.xinyadianwan.com/bbs/index.php"
	// var LoginPageUrl = "http://www.xinyadianwan.com/bbs/logging.php?action=login"
	var PS4PageURL = "http://www.xinyadianwan.com/bbs/exchangeps4/index.php?page_c=1&search_name=&gametype=&gamelang=&page="
	var targetGameNames = []string{"机战V", "深夜", "尼尔机械部队中文版", "进击的巨人", "奥丁领域中文版", "最终幻想15中文版", "心理测量者", "高达VS", "女神异闻录", "SD高达G世纪"}

	var notify *notificator.Notificator = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "新亚电玩",
	})
	index := 1

	// _ = getDocument(PS4PageURL + strconv.Itoa(1))
	// _ = getDocument(PS4PageURL + strconv.Itoa(15))
	// _ = getDocument(PS4PageURL + strconv.Itoa(20))

	for ; ; index++ {
		fmt.Println("请求第" + strconv.Itoa(index) + "页")
		doc := getDocument(PS4PageURL + strconv.Itoa(index))
		if doc == nil {
			continue
		}
		selections := doc.Find(".item")
		if selections == nil {
			break
		}
		if len(selections.Nodes) <= 0 {
			break
		}
		selections.Each(func(index int, item *goquery.Selection) {
			name := item.Find("h3").Text()
			// fmt.Println(name)
			for _, n := range targetGameNames {
				name = strings.TrimSpace(name)
				if strings.HasPrefix(name, n) {
					title := item.Find("h3").First().Text()
					storeNumber := item.Find(".price").First().Text()

					if storeNumber == "" {
						storeNumber = "暂无数据，可能是抽奖的"
					}

					fmt.Println("==============")
					fmt.Println(title)
					fmt.Println(storeNumber)
					go notify.Push(title, storeNumber, "", notificator.UR_NORMAL)
				}
			}
		})
		time.Sleep(1 * time.Second)
	}
}

func getDocument(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		fmt.Print("请求出错")
		fmt.Println(err.Error())
		panic(err)
	}
	defer res.Body.Close()
	// Convert the designated charset HTML to utf-8 encoded HTML.
	// `charset` being one of the charsets known by the iconv package.

	body, _ := ioutil.ReadAll(res.Body)
	rawHTMLString := string(body)
	convertedHTMLString, _ := iconv.ConvertString(rawHTMLString, "gbk", "utf-8")

	// use utfBody using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(convertedHTMLString))
	if err != nil {
		// handler error
		fmt.Print("解析出错")
		fmt.Println(err.Error())
		panic(err)
	}

	return doc
}
