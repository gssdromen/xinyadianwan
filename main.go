package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/gssdromen/notificator"
)

func main() {
	// var HomePageUrl = "http://www.xinyadianwan.com/bbs/index.php"
	// var LoginPageUrl = "http://www.xinyadianwan.com/bbs/logging.php?action=login"
	var PS4PageURL = "http://www.xinyadianwan.com/bbs/exchangeps4/index.php?page_c=1&search_name=&gametype=&gamelang=&page="

	var targetGameNames []string
	targetGameNames = []string{"深夜", "尼尔机械部队中文版", "进击的巨人", "最终幻想15中文版", "心理测量者", "高达VS", "女神异闻录5中", "SD高达G世纪", "丧尸围城", "怪物猎人", "洛克人遗产"}

	var notify *notificator.Notificator = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "新亚电玩",
	})
	index := 1
	hasGameAvailable := false

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
					store := item.Find(".price").First().Text()
					storeNumber := 0

					if store != "" {
						storeNumber, _ = strconv.Atoi(strings.Split(store, "：")[1])
					}

					fmt.Println("==============")
					fmt.Println(title)
					fmt.Println(storeNumber)
					if storeNumber > 0 {
						hasGameAvailable = true
						go notify.Push(title, store, 5, "default", PS4PageURL+strconv.Itoa(index), "", notificator.UR_NORMAL)
					}
				}
			}
		})
		time.Sleep(1 * time.Second)
	}
	// 如果一个可用的游戏都没有,显示一个提示
	if hasGameAvailable == false {
		notify.Push("很遗憾!", "", 3, "default", "", "", notificator.UR_NORMAL)
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

// func getTatgetGamesFormFile(filename string) (list []string) {
// 	file, _ := os.Open(filename)
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		line := strings.TrimSpace(line)
// 	}
// }
