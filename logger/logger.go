package logger

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func LogInsert(method, agent, path string, r *http.Request) {
	ip := getIP(r)
	f, _ := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	logstring := fmt.Sprintf("(%v)method: %v| ip: %v | User-Agent: %v| path: %v\n", time.Now().String(), method, ip, agent, path)
	f.WriteString(logstring)
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}
