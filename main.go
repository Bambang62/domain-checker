package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// Scrape WebsiteSEOchecker
func scrapeMoz(domain string) map[string]string {
	url := "https://websiteseochecker.com/bulk-check-page-authority/?query=" + domain
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RenderBot/1.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error WebsiteSEOchecker:", err)
		return nil
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	data := make(map[string]string)
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("td").First().Text())
		val := strings.TrimSpace(s.Find("td").Eq(1).Text())
		if key != "" {
			data[key] = val
		}
	})

	return data
}

// Scrape Ahrefs Free Backlink Checker
func scrapeAhrefs(domain string) map[string]string {
	url := "https://ahrefs.com/backlink-checker?input=" + domain
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RenderBot/1.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Ahrefs:", err)
		return nil
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	data := make(map[string]string)
	data["DR"] = strings.TrimSpace(doc.Find(".backlinks-profile-score").First().Text())
	data["Backlinks"] = strings.TrimSpace(doc.Find(".backlinks-profile-item:contains('Backlinks') .backlinks-profile-value").Text())
	data["Linking Websites"] = strings.TrimSpace(doc.Find(".backlinks-profile-item:contains('Linking websites') .backlinks-profile-value").Text())

	return data
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `
		<h2>🔍 Domain Quality Checker</h2>
		<form action="/check">
			<input name="domain" placeholder="example.com" style="width:250px"/>
			<button type="submit">Cek</button>
		</form>`)
	})

	r.GET("/check", func(c *gin.Context) {
		domain := c.Query("domain")
		if domain == "" {
			c.String(400, "Domain kosong")
			return
		}

		mozData := scrapeMoz(domain)
		ahrefsData := scrapeAhrefs(domain)
		wayback := "https://web.archive.org/web/*/" + domain

		html := "<h3>✅ Hasil Cek untuk: " + domain + "</h3><table border=1 cellpadding=5>"
		for k, v := range mozData {
			html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", k, v)
		}
		for k, v := range ahrefsData {
			html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", k, v)
		}
		html += fmt.Sprintf("<tr><td>Wayback Machine</td><td><a href='%s' target='_blank'>View History</a></td></tr>", wayback)
		html += "</table><br><a href='/'>🔙 Cek Domain Lain</a>"

		c.Header("Content-Type", "text/html")
		c.String(200, html)
	})

	// Ambil PORT dari Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default buat lokal
	}
	log.Println("Server jalan di port:", port)
	r.Run(":" + port)
}
