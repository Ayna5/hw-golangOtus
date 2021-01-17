package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var compile *regexp.Regexp

type keyValue struct {
	Key   string
	Value int
}

var sortedStruct []keyValue

func init() {
	pattern := ` ?(\S+.?) ?`
	compile = regexp.MustCompile(pattern)
}

func Top10(s string) (str []string) {
	res := make(map[string]int)
	if s == "" {
		return []string{}
	}

	words := compile.FindAllString(s, len(s))

	for _, w := range words {
		if w == "" {
			continue
		}

		w = strings.TrimFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		res[w]++
	}

	for key, value := range res {
		sortedStruct = append(sortedStruct, keyValue{key, value})
	}

	sort.Slice(sortedStruct, func(i, j int) bool {
		return sortedStruct[i].Value > sortedStruct[j].Value
	})
	fmt.Println(sortedStruct)

	i := 0
	for _, v := range sortedStruct {
		if i < 10 {
			str = append(str, v.Key)
		}
		i++
	}
	fmt.Println(str)
	return str
}
