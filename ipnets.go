package bslc

import (
    "net"
)

// IPNetContainer is a container for IP networks.
type IPNetContainer struct {
    networks []*net.IPNet
    cache map[string]bool
}

// Contains returns true if container includes the given host and false
// otherwise. The specified host can be a DNS name or an IP.
func (i *IPNetContainer) Contains(host string) bool {
    originalHost := host
    if result, exists := i.cache[host]; exists {
        return result
    }

    addToCache := func (result bool) bool {
        i.cache[originalHost] = result
        return result
    }

    ips, err := net.LookupHost(host)
    if err == nil {
        host = ips[0]
    }

    parsedIP := net.ParseIP(host)
    if parsedIP == nil {
        return addToCache(false)
    }

    for _, network := range i.networks {
        if network.Contains(parsedIP) {
            return addToCache(true)
        }
    }

    return addToCache(false)
}

// NewIPNetContainer returns a new IPNetContainer encompassing the networks specified.
// Networks are on CIDR notation form.
func NewIPNetContainer(networks []string) IPNetContainer {
    var parsedNetworks []*net.IPNet
    for _, network := range networks {
        _, ipnet, err := net.ParseCIDR(network)
        if err != nil { continue }
        parsedNetworks = append(parsedNetworks, ipnet)
    }

    return IPNetContainer{ networks: parsedNetworks, cache: make(map[string]bool) }
}
