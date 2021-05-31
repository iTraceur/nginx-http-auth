package utils

import (
	"testing"
)

func TestHttpAuth(t *testing.T) {
	ret := HttpAuth("demo", "123456")
	t.Log(ret)
}
