package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Handler utama
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// ambil domain dari query ?domain=
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		fmt.Fprintf(w, `
		<html>
			<head><title>Domain Checker</title></head>
			<body style="font-family:Arial,sans-serif;max-width:600px;margin:auto;">
				<h2>🌐 Domain Quality Checker</h2>
				<form method="GET">
					<input type="text" name="domain" placeholder="example.com" style="width:70%%;padding:8px;">
					<button type="submit" style="padding:8px;">Cek</button>
				</form>
			</body>
		</html>
		`)
		return
	}

	// --- TODO: disini tinggal tambahin logic scrape websiteseochecker + ahrefs + archive.org ---
	// untuk demo gua kasih mock hasil (bisa langsung ganti ke real scraping pake goquery)

	html := fmt.Sprintf(`
		<html>
			<head><title>Hasil Cek: %s</title></head>
			<body style="font-family:Arial,sans-serif;max-width:700px;margin:auto;">
				<h2>✅ Hasil Cek untuk: %s</h2>
				<table border="1" cellpadding="8" style="border-collapse:collapse;">
					<tr><td><b>DA (Moz)</b></td><td>31</td></tr>
					<tr><td><b>PA (Moz)</b></td><td>32</td></tr>
					<tr><td><b>Spam Score</b></td><td>1%%</td></tr>
					<tr><td><b>Total Backlinks</b></td><td>44K</td></tr>
					<tr><td><b>DoFollow</b></td><td>36%%</td></tr>
					<tr><td><b>NoFollow</b></td><td>64%%</td></tr>
					<tr><td><b>DR (Ahrefs)</b></td><td>20</td></tr>
					<tr><td><b>Linking Websites</b></td><td>330</td></tr>
					<tr><td><b>Domain Age</b></td><td>10 Tahun</td></tr>
					<tr><td><b>Wayback Machine</b></td>
						<td><a href="https://web.archive.org/web/*/%s" target="_blank">View History</a></td>
					</tr>
				</table>
				<br>
				<a href="/">🔙 Cek Domain Lain</a>
			</body>
		</html>
	`, domain, domain, domain)

	fmt.Fprint(w, html)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // biar bisa jalan di lokal juga
	}

	http.HandleFunc("/", handler)

	fmt.Println("🔥 Server jalan di port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
