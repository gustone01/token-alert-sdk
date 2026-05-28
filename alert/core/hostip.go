package core

import "net"

// ResolveHostIP 返回配置的 host IP；未配置时尝试探测本机 IPv4。
func ResolveHostIP(configured string) string {
	if configured != "" {
		return configured
	}
	return detectLocalIPv4()
}

func detectLocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if v4 := ipnet.IP.To4(); v4 != nil {
			return v4.String()
		}
	}
	return ""
}
