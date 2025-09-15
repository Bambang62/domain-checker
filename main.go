package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Halaman input
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Hasil cek
	r.POST("/check", func(c *gin.Context) {
		domain := c.PostForm("domain")
		if domain == "" {
			c.String(http.StatusBadRequest, "Domain tidak boleh kosong")
			return
		}

		// Scrape WebsiteSEOChecker
		mozDA, mozPA, spamScore, backlinks, doFollow, noFollow, domainAge := scrapeMoz(domain)

		// Scrape Ahrefs Free
		dr, linkingWeb := scrapeAhrefs(domain)

		// Wayback Machine
		wayback := fmt.Sprintf("https://web.archive.org/web/*/%s", domain)

		// Render hasil
		c.HTML(http.StatusOK, "result.html", gin.H{
			"Domain":     domain,
			"DA":         mozDA,
			"PA":         mozPA,
			"Spam":       spamScore,
			"Backlinks":  backlinks,
			"DoFollow":   doFollow,
			"NoFollow":   noFollow,
			"DR":         dr,
			"LinkingWeb": linkingWeb,
			"DomainAge":  domainAge,
			"Wayback":    wayback,
		})
	})

	log.Println("🚀 Server jalan di :8080")
	r.Run(":8080")
}

// ==============================
// Scraper Functions
// ==============================

func scrapeMoz(domain string) (string, string, string, string, string, string, string) {
	url := fmt.Sprintf("https://websiteseochecker.com/bulk-domain-authority-checker/?query=%s", domain)
	res, err := http.Get(url)
	if err != nil {
		return "-", "-", "-", "-", "-", "-", "-"
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	mozDA := doc.Find("td:contains('Moz Domain Authority')").Next().Text()
	mozPA := doc.Find("td:contains('Moz Page Authority')").Next().Text()
	spamScore := doc.Find("td:contains('Spam Score')").Next().Text()
	backlinks := doc.Find("td:contains('Total Backlinks')").Next().Text()
	doFollow := doc.Find("td:contains('DoFollow Backlinks')").Next().Text()
	noFollow := doc.Find("td:contains('NoFollow Backlinks')").Next().Text()
	domainAge := doc.Find("td:contains('Domain Age')").Next().Text()

	return mozDA, mozPA, spamScore, backlinks, doFollow, noFollow, domainAge
}

func scrapeAhrefs(domain string) (string, string) {
	url := fmt.Sprintf("https://ahrefs.com/backlink-checker?target=%s", domain)
	res, err := http.Get(url)
	if err != nil {
		return "-", "-"
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	dr := doc.Find(".domain-rating").Text()
	linkingWeb := doc.Find(".linking-websites").Text()

	return dr, linkingWeb
}