package iputils

import "testing"

func TestCheckIPInCIDR(t *testing.T) {
	ret:=CheckIPInCIDR("192.168.30.1","192.168.31.0/23")
	t.Log(ret)
}

func BenchmarkCheckIPInCIDR(b *testing.B) {
	ret:=CheckIPInCIDR("192.168.30.1","192.168.31.0/23")
	b.Log(ret)
}