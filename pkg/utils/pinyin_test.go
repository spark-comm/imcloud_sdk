package utils

import (
	"testing"
)

func Test_GetChineseFirstLetter(t *testing.T) {
	initial := GetChineseFirstLetter("掌声sdfsdfsd")
	t.Logf(initial)
}
