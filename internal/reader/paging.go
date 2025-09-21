package reader

import "strings"

func splitIntoPages(text string, minCharacterCount int) []string {
	paragraphs := strings.Split(text, "\r\n\r\n")

	pages := []string{}
	for _, p := range paragraphs {
		if len(pages) == 0 {
			pages = append(pages, "")
		}
		if len(pages[len(pages)-1]) > minCharacterCount {
			pages = append(pages, "")
		}
		if len(pages[len(pages)-1]) != 0 {
			pages[len(pages)-1] += "\r\n\r\n"
		}
		pages[len(pages)-1] += p
	}

	return pages
}
