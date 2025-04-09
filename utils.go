package deeplx_translator

import (
	"fmt"
	"slices"
	"strings"
)

var defaultSentenceTerminators = []string{
	".", ";", "!", "?", // English
	"。", "；", "！", "？", // CJK
	"\n", // Newline
}

func textToString(text any) (string, error) {
	switch v := text.(type) {
	case string:
		return v, nil
	case []string:
		return strings.Join(v, "\n"), nil
	default:
		return "", fmt.Errorf("unsupported text type")
	}
}

func textToStringSlice(text any) ([]string, error) {
	switch v := text.(type) {
	case string:
		if len(v) < 100 {
			return []string{v}, nil
		}
		return splitTextsAfter(v, defaultSentenceTerminators...), nil
	case []string:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported text type")
	}
}

func splitTextsAfter(text string, seps ...string) []string {
	results := []string{text}
	for _, sep := range seps {
		var temp []string
		for _, str := range results {
			parts := strings.SplitAfter(str, sep)
			temp = append(temp, parts...)
		}
		results = temp
	}
	return slices.DeleteFunc(results, func(s string) bool {
		return strings.TrimSpace(s) == ""
	})
}
