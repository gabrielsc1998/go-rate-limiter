package helpers

import (
	"net/http"
	"strings"
)

func GetIPFromRequest(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress != "" {
		return IPAddress
	}

	IPAddress = r.Header.Get("X-Forwarded-For")
	if IPAddress != "" {
		return IPAddress
	}

	IPAddress = r.RemoteAddr
	if IPAddress != "" {
		splittedIp := strings.Split(IPAddress, ":")
		return splittedIp[0]
	}

	return IPAddress
}
