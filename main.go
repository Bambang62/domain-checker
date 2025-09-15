package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"html/template"
)

type Result struct {
	Domain       string
	DA           string
	PA           string
	DR           string
	SpamScore    string
	Backlinks    string
	DoFollow     string
	NoFollow     string
	DomainAge    string
	ArchiveLink  string
}

func scrapeWebsiteSeoChecker(domain string) (da, pa, spam, backlinks, age string) {
	url := fmt.Sprintf("https://websiteseochecker.com/bulk-domain-authority-checker/?bulkurls=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return "-", "-", "-", "-", "-"
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// ambil value (harus disesuaikan selector html aslinya)
	da = strings.TrimSpace(doc.Find("td:contains('Moz Domain Authority')").Next().Text())
	pa = strings.TrimSpace(doc.Find("td:contains('Moz Page Authority')").Next().Text())
	spam = strings.TrimSpace(doc.Find("td:contains('Spam Score')").Next().Text())
	backlinks = strings.TrimSpace(doc.Find("td:contains('Total Backlinks')").Next().Text())
	age = strings.TrimSpace(doc.Find("td:contains('Domain Age')").Next().Text())

	if da == "" { da = "-" }
	if pa == "" { pa = "-" }
	if spam == "" { spam = "-" }
	if backlinks == "" { backlinks = "-" }
	if age == "" { age = "-" }

	return
}

func scrapeAhrefs(domain string) (dr, dofollow, nofollow string) {
	url := fmt.Sprintf("https://ahrefs.com/backlink-checker?input=%s", domain)

	resp, err := http.Get(url)
	if err != nil {
		return "-", "-", "-"
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	// ambil value (harus disesuaikan selector html aslinya)
	dr = strings.TrimSpace(doc.Find("div:contains('Domain Rating')").First().Text())
	dfPercent := strings.TrimSpace(doc.Find("div:contains('dofollow')").First().Text())
	nfPercent := strings.TrimSpace(doc.Find("div:contains('nofollow')").First().Text())

	if dr == "" { dr = "-" }
	if dfPercent == "" { dofollow = "-" } else { dofollow = dfPercent }
	if nfPercent == "" { nofollow = "-" } else { nofollow = nfPercent }

	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		domain := r.FormValue("domain")

		da, pa, spam, backlinks, age := scrapeWebsiteSeoChecker(domain)
		dr, df, nf := scrapeAhrefs(domain)

		data := Result{
			Domain:      domain,
			DA:          da,
			PA:          pa,
			DR:          dr,
			SpamScore:   spam,
			Backlinks:   backlinks,
			DoFollow:    df,
			NoFollow:    nf,
			DomainAge:   age,
			ArchiveLink: fmt.Sprintf("https://web.archive.org/web/*/%s", domain),
		}

		tmpl.Execute(w, data)
		return
	}
	tmpl.Execute(w, nil)
}

var tmpl = template.Must(template.New("ui").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Domain Quality Checker</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 40px; }
		form { margin-bottom: 20px; }
		table { border-collapse: collapse; width: 70%; }
		td, th { border: 1px solid #ddd; padding: 8px; }
		th { background: #f2f2f2; }
	</style>
</head>
<body>
	<h2>Domain Quality Checker</h2>
	<form method="POST">
		<input type="text" name="domain" placeholder="example.com" required>
		<button type="submit">Cek</button>
	</form>

	{{if .Domain}}
	<h3>Hasil Cek: {{.Domain}}</h3>
	<table>
		<tr><th>DA</th><td>{{.DA}}</td></tr>
		<tr><th>PA</th><td>{{.PA}}</td></tr>
		<tr><th>DR (Ahrefs)</th><td>{{.DR}}</td></tr>
		<tr><th>Spam Score</th><td>{{.SpamScore}}</td></tr>
		<tr><th>Total Backlinks</th><td>{{.Backlinks}}</td></tr>
		<tr><th>DoFollow</th><td>{{.DoFollow}}</td></tr>
		<tr><th>NoFollow</th><td>{{.NoFollow}}</td></tr>
		<tr><th>Domain Age</th><td>{{.DomainAge}}</td></tr>
		<tr><th>Archive.org</th><td><a href="{{.ArchiveLink}}" target="_blank">View History</a></td></tr>
	</table>
	{{end}}
</body>
</html>
`))

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server jalan di http://0.0.0.0:8080 🚀")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
