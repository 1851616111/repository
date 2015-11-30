package main

import (
	"testing"
)

func Test_buildTime(t *testing.T) {
	param := []string{"2013-01-02 15:04:05", "2014-01-02 15:04:05", "2015-05-26 15:04:05", "2015-10-28 15:04:05", "2015-11-28 15:04:05", "2015-11-30 15:04:05"}
	for i := 0; i < len(param); i++ {

		t.Log(param[i])
		t.Log(buildTime(param[i]))
	}
}
