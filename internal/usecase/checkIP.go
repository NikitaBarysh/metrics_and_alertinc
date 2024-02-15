package usecase

import (
	"fmt"
	"net"
	"net/http"
)

var IPCheckerVar *IPChecker

type IPChecker struct {
	subnet *net.IPNet
}

func (m *IPChecker) checkIP(clientIP string) bool {
	ip := net.ParseIP(clientIP)
	return m.subnet.Contains(ip)
}

func InitIPChecker(cidr string) (*IPChecker, error) {
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("err to init: %w", err)
	}
	return &IPChecker{subnet: subnet}, nil
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if IPCheckerVar != nil {
			ip := r.Header.Get("X-Real-IP")
			if ip == "" {
				fmt.Println("err to get header value")
				rw.WriteHeader(http.StatusForbidden)
				return
			}
			if !IPCheckerVar.checkIP(ip) {
				fmt.Println("wrong ip")
				rw.WriteHeader(http.StatusForbidden)
				return
			}
		}
		h.ServeHTTP(rw, r)

	})
}
