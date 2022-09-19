package common

import "strings"

func ContainInArray[T comparable](target T, targetArray []T) bool {
	for _, t := range targetArray {
		if target == t {
			return true
		}
	}
	return false
}

func ContainInArrayWithoutCapital(target string, targetArray []string) bool {
	for _, t := range targetArray {
		if strings.EqualFold(target, t) {
			return true
		}
	}
	return false
}
