package main

import (
	"regexp"
	"strings"
)

// Simple HTML to plain text converter
// This is a very basic implementation
func StripHTML(html string) string {
	// Remove HTML tags
	re := regexp.MustCompile("<[^>]*>")
	text := re.ReplaceAllString(html, "")
	
	// Replace common entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	
	// Replace consecutive whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	
	// Remove leading/trailing whitespace
	return strings.TrimSpace(text)
}

//WARNING, THIS WILL NOT BE IN PRODUCTION, HTML RENDERING WILL BE FULLY MADE BEFORE THIS APP IS FINISHED.