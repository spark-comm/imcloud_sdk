package utils

import (
	"testing"
)

func Test_GetChineseFirstLetter(t *testing.T) {
	initial := GetChineseFirstLetter("惠")
	t.Logf(initial)
}
