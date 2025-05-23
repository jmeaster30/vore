package algo

import (
	"strings"
)

func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x int, y int) int {
	if x > y {
		return y
	}
	return x
}

func Window[T any](array []T, size int) [][]T {
	var results [][]T
	for idx := range array {
		results = append(results, array[idx:Min(idx+size, len(array))])
	}
	return results
}

func SplitKeep(target string, split string) []string {
	var results []string
	for len(target) != 0 {
		splitStart := strings.Index(target, split)
		if splitStart == 0 {
			target = target[len(split):]
			results = append(results, split)
		} else if splitStart == -1 {
			results = append(results, target)
			break
		} else {
			substr := target[0:splitStart]
			target = target[splitStart:]
			results = append(results, substr)
		}
	}
	return results
}

func Combine[T any](array ...T) []T {
	return array
}
