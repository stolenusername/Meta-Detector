package main

import (
	"fmt"
	"os"
	"strings"
)

func generateSearchURLs(domain string) []string {
	formats := []string{"pdf", "doc", "docx", "dot", "odt", "rtf", "xls", "xlsx", "ods", "ppt", "pptx", "pps", "txt"}

	// Additional search strings
	searchStrings := []string{"Index of", ".yaml", ".htaccess", ".circleci", ".git", ".ssh", "id_rsa", "config.inc.php"}

	var urls []string
	for _, format := range formats {
		for _, searchStr := range searchStrings {
			url := fmt.Sprintf("https://www.google.com/search?q=site:%s+filetype:%s+%s", domain, format, searchStr)
			urls = append(urls, url)
		}
	}

	return urls
}

func generateHTMLPage(domain string, urls []string) string {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString("<html>\n")
	sb.WriteString("<head>\n")
	sb.WriteString("<title>OSINT Search Results for ")
	sb.WriteString(domain)
	sb.WriteString("</title>\n")
	sb.WriteString("</head>\n")
	sb.WriteString("<body>\n")
	sb.WriteString("<h1>OSINT Search Results for ")
	sb.WriteString(domain)
	sb.WriteString("</h1>\n")
	sb.WriteString("<ul>\n")
	for _, url := range urls {
		sb.WriteString("<li><a href=\"")
		sb.WriteString(url)
		sb.WriteString("\" target=\"_blank\">")
		sb.WriteString(url)
		sb.WriteString("</a></li>\n")
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
