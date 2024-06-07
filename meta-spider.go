package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Configuration struct to hold spider configuration
type Config struct {
	Domain         string
	Depth          int
	Concurrency    int
	SearchItems    []string
	TitleWords     []string
	FileExtensions []string
	URLWords       []string
}

// Spider struct to hold spider state
type Spider struct {
	Config
	Visited sync.Map
	Mutex   sync.Mutex
}

func main() {
	// Parse command line arguments
	domain := flag.String("domain", "", "Domain to spider")
	depth := flag.Int("depth", 1, "Depth to spider")
	concurrency := flag.Int("concurrency", 10, "Concurrency level")
	configFile := flag.String("config", "spider.config", "Spider configuration file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// If no arguments are provided, print usage information
	if *domain == "" {
		flag.Usage()
		return
	}

	// Read spider configuration from file
	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Println("Error reading spider configuration:", err)
		return
	}

	config.Domain = *domain
	config.Depth = *depth
	config.Concurrency = *concurrency

	spider := Spider{
		Config: config,
	}

	spider.crawl(*domain, *depth)
}

// readConfig reads spider configuration from file
func readConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Skipping invalid line in configuration file: %s\n", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "search_words":
			words := strings.Split(value, ",")
			for _, w := range words {
				config.SearchItems = append(config.SearchItems, strings.TrimSpace(w))
			}
		case "title_words":
			words := strings.Split(value, ",")
			for _, w := range words {
				config.TitleWords = append(config.TitleWords, strings.TrimSpace(w))
			}
		case "file_extensions":
			extensions := strings.Split(value, ",")
			for _, ext := range extensions {
				config.FileExtensions = append(config.FileExtensions, strings.TrimSpace(ext))
			}
		case "url_words":
			words := strings.Split(value, ",")
			for _, w := range words {
				config.URLWords = append(config.URLWords, strings.TrimSpace(w))
			}
		}
	}
	return config, scanner.Err()
}

// crawl starts the spidering process
func (s *Spider) crawl(domain string, depth int) {
	var wg sync.WaitGroup
	queue := make(chan string)

	wg.Add(1)
	go func() {
		queue <- domain
		wg.Done()
	}()

	for i := 0; i < s.Config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for link := range queue {
				s.visit(link, depth, queue)
			}
		}()
	}

	wg.Wait()
	close(queue)
}

// visit retrieves the HTML content of a URL and searches for specified items
func (s *Spider) visit(link string, depth int, queue chan string) {
	if depth <= 0 {
		return
	}

	// Check if the link is within the specified domain
	if !strings.Contains(link, s.Domain) {
		return
	}

	// Check if the link has already been visited
	if _, visited := s.Visited.LoadOrStore(link, true); visited {
		return
	}

	fmt.Println("Visiting:", link)

	// Check if the link has a scheme, if not, prepend "http://"
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "http://" + link
	}

	resp, err := http.Get(link)
	if err != nil {
		// If http request failed, try with https
		link = strings.Replace(link, "http://", "https://", 1)
		resp, err = http.Get(link)
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return
		}
	}
	defer resp.Body.Close()

	baseURL, err := url.Parse(link)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	base := baseURL.Scheme + "://" + baseURL.Host

	if baseURL.Path == "" {
		base += "/"
	}

	// Read the response body as a string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	bodyString := string(bodyBytes)

	tokenizer := html.NewTokenizer(strings.NewReader(bodyString))

	// Check if the URL points to a file with the specified extensions
	for _, ext := range s.FileExtensions {
		if strings.HasSuffix(link, ext) {
			fmt.Printf("Found file with extension '%s' at %s\n", ext, link)
			// Log the found file to a file
			s.logToFile("File with extension "+ext, link)
		}
	}

	// Search for words in the title of the page
	titleWords := strings.Join(s.TitleWords, "|")
	titleRegex := regexp.MustCompile(`(?i)<title.*?>(.*?)</title>`)
	titleMatch := titleRegex.FindStringSubmatch(bodyString)
	if len(titleMatch) > 1 {
		title := titleMatch[1]
		if strings.Contains(strings.ToLower(title), strings.ToLower(titleWords)) {
			fmt.Printf("Found title containing '%s' in %s\n", titleWords, link)
			s.logToFile(fmt.Sprintf("Title containing '%s'", titleWords), link)
		}
	}

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						url := attr.Val
						if strings.HasPrefix(url, "/") {
							url = base + url
						}
						if strings.Contains(url, s.Domain) {
							go func() {
								queue <- url
							}()
						}
						break
					}
				}
			}
			// Search for other items here
			for _, searchItem := range s.SearchItems {
				if strings.Contains(strings.ToLower(token.String()), strings.ToLower(searchItem)) {
					fmt.Printf("Found '%s' in %s\n", searchItem, link)
					// Log the found item to a file
					s.logToFile(searchItem, link)
				}
			}
		}
	}

	// Check if the URL contains specified words
	for _, urlWord := range s.URLWords {
		if strings.Contains(link, urlWord) {
			fmt.Printf("Found URL containing '%s' at %s\n", urlWord, link)
			// Log the found URL to a file
			s.logToFile("URL containing "+urlWord, link)
		}
	}
}

// logToFile appends the found item and the link to a file
func (s *Spider) logToFile(searchItem, link string) {
	fileName := fmt.Sprintf("%s-crawl.txt", s.Domain)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("Found '%s' in %s\n", searchItem, link)); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
