package main


import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

)

func main () {
	c := colly.NewCollector(
		colly.AllowedDomains("emojipedia.org"),
	)
	//callback from scraped page which contains an article element
	c.OnHTML("article", func(e *colly.HTMLElement) {
		isEmojiPage := false

		//get the meta tags from page
		metaTags := e.DOM.ParentsUntil("~").Find("meta")
		metaTags.Each(func(_ int, s *goquery.Selection) {
			//Search for og:type meta tags
			property, _ := s.Attr("property")
			if strings.EqualFold(property, "og:type") {
				content, _ := s.Attr("content")

				//Emojie pages have "article" as their ogtype
				isEmojiPage = strings.EqualFold(content, "article")
			}
		})

		if isEmojiPage {
			//locate the emoji title page
			fmt.Println("Emoji: ", e.DOM.Find("h1").Text())
			//Grab all the text from the emojie description
			fmt.Println("Description: ", e.DOM.Find(".description").Find("p").Text())
		}
	})

	//Callback for links on scraped pages
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//Get the linked URL from the anchor tag
		link := e.Attr("href")
		//webcrawler visits the linked URL
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly .Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://emojipedia.org")
}