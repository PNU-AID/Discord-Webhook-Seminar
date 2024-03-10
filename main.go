package main

import (
	"log"
	"os"
	"strings"

	"github.com/disgoorg/disgo/webhook"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

// https://1minute-before6pm.tistory.com/42

func main() {
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .eng file")
	}
	discord_url := os.Getenv("DISCORD_URL")

	client, webhook_err := webhook.NewWithURL(discord_url)
	if webhook_err != nil {
		log.Fatal("Error in discord_url")
	}

	// now := time.Now()
	web_url := "https://dev-event.vercel.app/events"

	var data = []map[string]string{}
	collector := colly.NewCollector()
	collector.OnHTML(".list_wrapper__tpe4x", func(e *colly.HTMLElement) {
		var msg = map[string]string{}

		e.ForEach(".Item_item__86e_I", func(i int, e_1 *colly.HTMLElement) {
			var flag bool = false
			if e_1.ChildText(".DdayTag_tag__YDKTM") == "Today" {
				flag = true
			}
			title := e_1.ChildText(".Item_item__content__title___fPQa")
			msg["title"] = title
			tag := e_1.ChildText(".FilterTag_tag__JRNsk")
			is_contain := strings.Contains(tag, "AI")
			if is_contain {
				url := e.ChildAttr("a", "href")
				msg["desc"] = url + "\n"
				e_1.ForEach(".Item_wrap__qzpSH", func(i int, e_2 *colly.HTMLElement) {
					desc_title := e_2.ChildText(".Item_label__vgqgX")
					desc_content_1 := e_2.ChildText(".Item_host__zNXMy")
					desc_content_2 := e_2.ChildText(".Item_date__kVMJZ")
					msg["desc"] += desc_title + " " + desc_content_1 + desc_content_2 + "\n"
				})
				if flag {
					data = append(data, msg)
				}
			}
		})
	})

	collector.Visit(web_url)

	var content string = "## 진행중인 개발자 행사\n"
	for _, d := range data {
		title := d["title"]
		desc := d["desc"]
		content += "### " + title + "\n" + desc + "\n"
	}

	_, err := client.CreateContent(content)
	if err != nil {
		log.Fatal(err)
	}

}
