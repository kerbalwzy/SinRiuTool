package utils

import "testing"

func TestIsNumber(t *testing.T) {
	value := "3.14"
	res, yes := IsNumber(value)
	t.Log(yes)
	t.Log(res)
}
