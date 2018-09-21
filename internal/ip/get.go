package ip

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func GetIP() (net.IP, error) {
	url := "https://icanhazip.com"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	str := strings.TrimSpace(string(content))
	ip := net.ParseIP(str)
	if ip == nil {
		return nil, errors.New("invalid IP returned")
	}

	return ip, nil
}
