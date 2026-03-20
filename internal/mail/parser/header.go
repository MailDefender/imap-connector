package parser

import (
	"fmt"
	"slices"
	"strings"
)

type headerParser struct {
	rawContent    string
	parsedContent map[string][]string
}

var patternToReplace = map[string]string{
	"\r": "",
	"\t": " ",
	"  ": " ",
}

func NewHeaderParser(input string) headerParser {
	return headerParser{
		rawContent:    input,
		parsedContent: make(map[string][]string),
	}
}

func cleanHeader(value string) string {
	out := value

	for pattern, replacement := range patternToReplace {
		out = strings.ReplaceAll(out, pattern, replacement)
	}

	// Manage the case if we have 3 spaces in a row which was replaced by 2 spaces in the previous loop
	out = strings.ReplaceAll(out, "  ", " ")
	out = strings.TrimSpace(out)

	return out
}

func removePrefix(in string, prefix []string) string {
	if len(in) == 0 {
		return ""
	}

	if slices.Contains(prefix, string(in[0])) {
		return removePrefix(in[1:], prefix)
	}
	return in
}

func (hp *headerParser) Parse() error {
	lines := strings.Split(hp.rawContent, "\n")

	var key, value string

	for _, line := range lines {
		if line == "" || line == "\r" || line == "\n" {
			continue
		}

		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
			value = fmt.Sprintf("%s %s", value, removePrefix(line, []string{" "}))
		} else {
			if key != "" && value != "" {
				if _, exist := hp.parsedContent[key]; exist {
					hp.parsedContent[key] = append(hp.parsedContent[key], cleanHeader(value))
				} else {
					hp.parsedContent[key] = []string{cleanHeader(value)}
				}
			}

			splittedLine := strings.Split(line, ":")

			if len(splittedLine) < 2 {
				fmt.Printf("No enought items : %s", line)
			}

			key = splittedLine[0]
			value = strings.Join(splittedLine[1:], ":")
			if strings.HasPrefix(value, " ") {
				value = value[1:]
			}
		}
	}
	if key != "" && value != "" {
		if _, exist := hp.parsedContent[key]; exist {
			hp.parsedContent[key] = append(hp.parsedContent[key], cleanHeader(value))
		} else {
			hp.parsedContent[key] = []string{cleanHeader(value)}
		}
	}

	return nil
}

func (hp headerParser) Get(key string) []string {
	return hp.parsedContent[key]
}

func (hp headerParser) GetAll() map[string][]string {
	return hp.parsedContent
}
