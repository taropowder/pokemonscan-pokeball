package utils

import (
	"fmt"
	"testing"
)

func TestGetCpuPercent(t *testing.T) {
	fmt.Println(GetCpuPercent())
	fmt.Println(GetMemPercent())
}
