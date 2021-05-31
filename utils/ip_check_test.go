package utils

import (
	"testing"
)

func TestIpCheck(t *testing.T) {
	checkIps := []string{
		"10.10.10.255",
		"10.10.11.254",
		"10.10.12.25",
		"10.10.10.256",
	}
	ips := []string{
		"10.10.10.0/24",
		"10.10.11.1-10.10.11.254",
		"10.10.12.25",
	}
	for _, checkIp := range checkIps {
		if IpCheck(checkIp, ips) {
			t.Log(true)
		} else {
			t.Log(false)
		}
	}
}

func TestInSlice(t *testing.T) {
	sss := []string{"1", "2", "3", "4", "5", "6"}
	t.Log(InSlice("1", sss))
}
