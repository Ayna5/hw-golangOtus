package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

type keyValue struct {
	Key   string
	Value int
}

var sortedStruct []keyValue

func Top10(s string) (str []string) {
	if s == "" {
		return []string{}
	}
	res := make(map[string]int)

	ww := strings.Fields(s)
	for _, val := range ww {
		w := strings.ToLower(val)
		if w != "" && w != "-" {
			res[w]++
		}
	}

	for key, value := range res {
		sortedStruct = append(sortedStruct, keyValue{key, value})
	}

	sort.Slice(sortedStruct, func(i, j int) bool {
		return sortedStruct[i].Value > sortedStruct[j].Value
	})

	if len(sortedStruct) > 10 {
		return words(sortedStruct[:10])
	}
	return words(sortedStruct)
}

func words(arr []keyValue) []string {
	str := make([]string, 0, len(arr))
	for _, v := range arr {
		str = append(str, v.Key)
	}
	return str
}
