# Google Dork Search Generator

This program generates an HTML file containing Google Dork search queries based on the parameters specified in the `search.config` file. It allows users to easily perform Google Dork searches for various purposes, such as finding specific file types, directories, login pages, and sensitive information.

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
2. Run the program with the domain as an argument:
3. The program will generate an HTML file named `<domain>_search_results.html` containing the Google Dork search results for the specified domain.

### Example Output

For example, if the domain is `domain.com`, the program will generate an HTML file named `domain.com_search_results.html` containing the Google Dork search results for `domain.com`.

