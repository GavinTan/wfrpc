package tools

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IndexOf(arr []string, e string) int {
	for k, v := range arr {
		if e == v {
			return k
		}
	}
	return -1
}

func RandomPort() string {
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(20000-10000) + 10000
	return strconv.Itoa(port)
}

func GetHostIP() (ip string) {
	defer func() {
		if r := recover(); r != nil {
			ip = "127.0.0.1"
		}
	}()

	conn, _ := net.Dial("udp", "8.8.8.8:53")
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]

	return ip
}
