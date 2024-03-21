package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

// readSearchConfig reads search parameters and descriptions from search.config file
func readSearchConfig(filename string) (map[string]string, error) {
	config := make(map[string]string)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid search parameter format: %s", line)
		}
		config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return config, nil
}

// generateSearchURLs generates Google search URLs for the given domain based on search parameters
func generateSearchURLs(domain string, searches map[string]string) []string {
	var urls []string
	for search, desc := range searches {
		encodedSearch := url.QueryEscape(search + " site:" + domain)
		url := fmt.Sprintf("https://www.google.com/search?q=%s", encodedSearch)
		urls = append(urls, fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, desc))
	}
	return urls
}

// generateHTMLPage generates HTML page with search results
func generateHTMLPage(domain string, urls []string) string {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString("<html>\n")
	sb.WriteString("<head>\n")
	sb.WriteString("<title>Meta-Detector Results for ")
	sb.WriteString(domain)
	sb.WriteString("</title>\n")
	sb.WriteString("</head>\n")
	sb.WriteString("<body>\n")
	sb.WriteString("<h1>Meta-Detector Results for ")
	sb.WriteString(domain)
	sb.WriteString("</h1>\n")
	sb.WriteString("<ul>\n")
	for _, url := range urls {
		sb.WriteString("<li>")
		sb.WriteString(url)
		sb.WriteString("</li>\n")
	}
	sb.WriteString("</ul>\n")
	sb.WriteString("<button onclick=\"openAllLinks()\">Open All Links in New Tab</button>\n")
	sb.WriteString("<script>\n")
	sb.WriteString("function openAllLinks() {\n")
	sb.WriteString("    var links = document.getElementsByTagName('a');\n")
	sb.WriteString("    var delay = 1000; // 1 second delay\n")
	sb.WriteString("    for (var i = 0; i < links.length; i++) {\n")
	sb.WriteString("        setTimeout(function(link) {\n")
	sb.WriteString("            window.open(link.href, '_blank');\n")
	sb.WriteString("        }, delay * i, links[i]);\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	sb.WriteString("</script>\n")
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

	// Read search parameters from search.config
	searches, err := readSearchConfig("search.config")
	if err != nil {
		fmt.Println("Error reading search config:", err)
		return
	}

	// Generate search URLs based on search parameters
	urls := generateSearchURLs(domain, searches)

	// Generate HTML page with search results
	html := generateHTMLPage(domain, urls)

	// Write HTML page to file
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
