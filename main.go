package main

import (
	"context"
	"fmt" //ms test
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/disgoorg/disgo/webhook" //Discord 관련
	"github.com/joho/godotenv"
)

// https://1minute-before6pm.tistory.com/42

func main() {
	env_err := godotenv.Load() //Discord 관련
	if env_err != nil {
		log.Fatal("Error loading .eng file")
	}
	discord_url := os.Getenv("DISCORD_URL")

	client, webhook_err := webhook.NewWithURL(discord_url) //Discord 관련
	if webhook_err != nil {
		log.Fatal("Error in discord_url")
	}

	// now := time.Now() //Discord 관련
	web_url := "https://dev-event.vercel.app/events"

	var data []map[string]string

	ctx, cancel := chromedp.NewContext(context.Background()) //, chromedp.WithDebugf(log.Printf)) // to debugging
	defer cancel()

	// Visit the webpage and extract the relevant information
	var innerNodes []*cdp.Node
	err_4 := chromedp.Run(ctx,
		chromedp.Navigate(web_url),
		chromedp.Sleep(3*time.Second), // 페이지 로딩을 위한 시간 대기
		chromedp.WaitVisible(".Home_section__EaDnq", chromedp.ByQuery),
		chromedp.Nodes(".Item_item__container___T09W", &innerNodes, chromedp.ByQueryAll),
	)
	if err_4 != nil {
		log.Fatal("Error HTML", err_4)
	}

	for _, node := range innerNodes {
		var todayText string
		err_1 := chromedp.Run(ctx,
			chromedp.Text(".DdayTag_tag__6_oE7", &todayText, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err_1 != nil {
			log.Println("Error extracting Today flag:", err_1)
			continue
		} else if !strings.Contains(todayText, "Today") {
			continue
		}

		var tagNodes []*cdp.Node
		err_2 := chromedp.Run(ctx,
			chromedp.Nodes(".FilterTag_tag__etNfv", &tagNodes, chromedp.ByQueryAll, chromedp.FromNode(node)),
		)
		if err_2 != nil {
			log.Println("Error extracting Tag:", err_2)
			continue
		}
		isAI := false
		for _, tagNode := range tagNodes {
			// fmt.Println(tagNode.Children[0].NodeValue)
			if strings.Contains(tagNode.Children[0].NodeValue, "AI") {
				isAI = true
				break
			}
		}
		fmt.Println()
		if !isAI {
			continue
		}

		var msg = map[string]string{}
		var titleText, urlText, dateText, hostText string
		err_3 := chromedp.Run(ctx,
			chromedp.Text(".Item_item__content__title__94_8Q", &titleText, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("a", "href", &urlText, nil, chromedp.FromNode(node)),
			chromedp.Text(".Item_date__date__CoMqV", &dateText, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".Item_host__3dy8_", &hostText, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err_3 != nil {
			log.Println("Error extracting title or url:", err_3)
			continue
		}
		msg["title"] = titleText
		msg["desc"] = urlText + "\n" + "주최: " + hostText + "\n" + "모집: " + dateText

		data = append(data, msg)
	}

	var content string = "## 진행중인 개발자 행사\n"
	for _, d := range data {
		content += "###" + d["title"] + "\n"
		content += d["desc"] + "\n\n"
	}

	fmt.Println("Webhook content preview:\n", content) //ms test
	fmt.Println("노드 개수: ", len(innerNodes))            //ms test
	fmt.Println("필터링된 노드 개수: ", len(data))             //ms test

	_, err_4 = client.CreateContent(content)
	if err_4 != nil {
		log.Fatal(err_4)
	}
}
