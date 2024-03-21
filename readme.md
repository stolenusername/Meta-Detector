# Meta-Detector
![Meta Detector](images/meta-detector.jpg)
<h3 style="text-align: center">


[![License: MIT](https://img.shields.io/badge/License-MIT-darkred.svg)](https://github.com/stolenusername/Meta-Detector/blob/main/LICENSE)
[![made-with-Go](https://img.shields.io/badge/Made%20with-GoLang-blue.svg)](https://go.dev/)

</h3>

This program generates an HTML file containing Google Dork search queries based on the parameters specified in the `search.config` file. It allows users to easily perform Google Dork searches for various purposes, such as finding specific file types, directories, login pages, and sensitive information.

## Downloading
1. Create the desired directory on the sytem you will install Meta-Detector on.
2. From the command line within that directory: `git clone https://github.com/stolenusername/Meta-Detector/`.
3. `cd Meta-Detector`
4. `go build meta-detector.go` (This step requires that you have Go installed on the system).

## How It Works

1. The program reads search parameters and descriptions from the `search.config` file.
2. It generates Google search URLs for the given domain based on the search parameters.
3. The URLs are embedded in an HTML page, which is then written to a file.

## `search.config` File

The `search.config` file contains a list of search parameters and their descriptions. Each line in the file follows the format:
- `<search_parameter>`: Specifies the Google Dork search parameter.
- `<Description>`: Provides a brief description of what the search parameter targets.

The search parameters can include operators such as `site:`, `filetype:`, `inurl:`, `intitle:`, `intext:`, and logical operators like `OR`.

### Updating `search.config`

To update the `search.config` file:
1. Open the `search.config` file in a text editor.
2. Add or modify the search parameters following the specified format (`<search_parameter> | <Description>`).
3. Save the changes to the `search.config` file.

## Usage

1. Ensure that the `search.config` file is correctly configured with the desired search parameters.
2. Run the program with the domain as an argument: `./meta-detector domain.com`
3. The program will generate an HTML file named `<domain>_search_results.html` containing the Google Dork search results for the specified domain.
4. Each link in the HTML file opens the Google search results in a new tab.
5. There is a button at the bottom of the links to open all the links in new tabs at once.

### Example Output

For example, if the domain is `domain.com`, the program will generate an HTML file named `domain.com_search_results.html` containing the Google Dork search results for `domain.com`.

## Troubleshooting

When using Meta-Detector tool, there are a few things to keep in mind:

### Pop-up Blocker:
Your browser's pop-up blocker might prevent the links from opening in new tabs. You may need to disable it or allow pop-ups for the current page to ensure that all links open successfully.

### "Prove You Are Not a Robot" Prompt:
Due to the nature of how Google throttles Google Dork searches, you may encounter the "Prove You Are Not a Robot" prompt, especially if you're performing a large number of searches in a short period. If you encounter this prompt, you'll need to complete the CAPTCHA verification to continue.

### Adjusting Delay:
In the code, there's a delay set before each link is opened in a new tab. If you find that you are getting throttled by Google, you can adjust this delay. To do so, locate the following line in the code: `var delay = 1000; // 1 second delay`
