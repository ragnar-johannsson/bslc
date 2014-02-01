package bslc

import "testing"

func TestIPNetContainer(t *testing.T) {
    subnets := []string {
        "1.2.3.4/24",
        "4.3.2.1/16",
        "8.8.8.8/32",
        "invalid/",
    }

    ipNets := NewIPNetContainer(subnets)

    if !ipNets.Contains("1.2.3.255") {
        t.Error("ipNets does not contain 1.2.3.255")
    }

    if ipNets.Contains("100.100.100.100") {
        t.Error("ipNets contains 100.100.100.100")
    }

    if !ipNets.Contains("google-public-dns-a.google.com") {
        t.Error("ipNets does not contain google public dns server")
    }

    if ipNets.Contains("doesnotresolve") {
        t.Error("ipNets contains unresolvable host")
    }

    ipNets.Contains("1.2.3.255")
    if len(ipNets.cache) != 4 {
        t.Error("ipNets.cache length != 4")
    }
}

