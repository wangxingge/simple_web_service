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

func WebHandle_Add(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v web add request.", hostname, ip)))
}
func WebHandle_Remove(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v web remove request.", hostname, ip)))
}
func WebHandle_Search(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v web search request.", hostname, ip)))
}

func PCHandle_Add(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v pc add request.", hostname, ip)))
}
func PCHandle_Remove(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v pc remove request.", hostname, ip)))
}
func PCHandle_Search(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v pc search request.", hostname, ip)))
}

func MobileHandle_Add(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v mobile add request.", hostname, ip)))
}
func MobileHandle_Remove(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v mobile remove request.", hostname, ip)))
}
func MobileHandle_Search(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v mobile search request.", hostname, ip)))
}

var (
	router = httprouter.New()
)

func main() {

	hostname, _ = os.Hostname()
	ip = getInternalIP()
	router.GET("/ping", Ping)
	router.GET("/who", Whoami)
	router.GET("/web_root/web/add", WebHandle_Add)
	router.GET("/web_root/web/remove", WebHandle_Remove)
	router.GET("/web_root/web/search", WebHandle_Search)
	router.GET("/web_root/mobile/add", MobileHandle_Add)
	router.GET("/web_root/mobile/remove", MobileHandle_Remove)
	router.GET("/web_root/mobile/search", MobileHandle_Search)
	router.GET("/web_root/pc/add", PCHandle_Add)
	router.GET("/web_root/pc/remove", PCHandle_Remove)
	router.GET("/web_root/pc/search", PCHandle_Search)
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
