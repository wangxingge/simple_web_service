package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	hostname = ""
	ip       = ""
)

func Ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte("ok"))
}

func Whoami(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("I am %v, address %v", hostname, ip)))
}

func WebHandle(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v web request.", hostname, ip)))
}

func PCHandle(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v pc request.", hostname, ip)))
}

func MobileHandle(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v mobile request.", hostname, ip)))
}

var (
	router = httprouter.New()
)

func main() {
	hostname, _ = os.Hostname()
	ip = getInternalIP()
	router.GET("/ping", Ping)
	router.GET("/who", Whoami)
	router.GET("/web_root/web", WebHandle)
	router.GET("/web_root/mobile", MobileHandle)
	router.GET("/web_root/pc", PCHandle)
	log.Fatal(http.ListenAndServe("0.0.0.0:80", router))
}

func getInternalIP() string {

	addressList, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var ipAddr string
	for _, a := range addressList {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {

				ipAddr = ipNet.IP.String()
			}
		}
	}
	return ipAddr
}
