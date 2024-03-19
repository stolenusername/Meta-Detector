package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func generateSearchURLs(domain string) []string {
	dorkSearches := map[string]string{
		"site:" + domain:                                "Site-specific search for domain",
		"site:*.example.com":                            "Site-specific search for subdomains",
		"filetype:pdf site:" + domain:                  "PDF files on domain",
		"filetype:doc site:" + domain:                  "Word documents on domain",
		"filetype:xls site:" + domain:                  "Excel spreadsheets on domain",
		"filetype:txt site:" + domain:                  "Text files on domain",
		"inurl:login.php site:" + domain:               "Login pages on domain",
		"intitle:\"login page\" site:" + domain:        "Pages with 'login page' in title on domain",
		"intext:\"username\" intext:\"password\" site:" + domain: "Pages containing 'username' and 'password' on domain",
		"inurl:intitle:\"index of /\" site:" + domain:        "Directories listing files on domain",
		"inurl:intitle:\"index of /backup\" site:" + domain:  "Backup directories on domain",
		"inurl:intitle:\"index of /config\" site:" + domain:  "Configuration files on domain",
		"inurl:intitle:\"phpinfo()\" site:" + domain:         "PHP configuration details on domain",
		"inurl:intitle:\"Welcome to phpMyAdmin\" site:" + domain: "phpMyAdmin installations on domain",
		"inurl:intitle:\"Welcome to OpenSSH\" site:" + domain: "OpenSSH installations on domain",
		"inurl:intitle:\"Error Occurred\" OR intitle:\"Server Error\" site:" + domain: "Pages with error messages on domain",
		"inurl:intext:\"MySQL error\" site:" + domain:        "MySQL error messages on domain",
		"inurl:intitle:\"Welcome to nginx\" site:" + domain:  "nginx web server installations on domain",
		"inurl:intitle:\"Apache2 Ubuntu Default Page\" site:" + domain: "Default Apache pages on Ubuntu on domain",
		"inurl:intitle:\"Index of\" intext:\"Served by Serv-U\" site:" + domain: "Serv-U FTP servers on domain",
		"inurl:intitle:\"Webcam Live Image\" site:" + domain: "Webcams broadcasting live images on domain",
		"inurl:intitle:\"Network Camera NetworkCamera\" site:" + domain: "Network cameras on domain",
		"inurl:intitle:\"Live View / - AXIS\" site:" + domain: "AXIS network cameras on domain",
	}

	var urls []string
	for search, desc := range dorkSearches {
		encodedSearch := url.QueryEscape(search)
		url := fmt.Sprintf("https://www.google.com/search?q=%s", encodedSearch)
		urls = append(urls, fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, desc))
	}
	return urls
}

func generateHTMLPage(domain string, urls []string) string {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString("<html>\n")
	sb.WriteString("<head>\n")
	sb.WriteString("<title>Google Dork Search Results for ")
	sb.WriteString(domain)
	sb.WriteString("</title>\n")
	sb.WriteString("</head>\n")
	sb.WriteString("<body>\n")
	sb.WriteString("<h1>Google Dork Search Results for ")
	sb.WriteString(domain)
	sb.WriteString("</h1>\n")
	sb.WriteString("<ul>\n")
	for _, url := range urls {
		sb.WriteString("<li>")
		sb.WriteString(url)
		sb.WriteString("</li>\n")
	}
	sb.WriteString("</ul>\n")
	sb.WriteString("</body>\n")
	sb.WriteString("</html>")

	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <domain>")
		return
	}

	domain := os.Args[1]
	urls := generateSearchURLs(domain)
	html := generateHTMLPage(domain, urls)

	file, err := os.Create(domain + "_search_results.html")
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		fmt.Println("Error writing to HTML file:", err)
		return
	}

	fmt.Println("HTML page generated successfully:", domain+"_search_results.html")
}
