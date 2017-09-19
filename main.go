package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
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

func Deep(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v deep request.", hostname, ip)))
}
func Deep_Sub(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v deep_sub request.", hostname, ip)))
}

func Restful(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("%v handle client %v restful category:%v, index:%v, id:%v, mark:%v request.",
		hostname, ip, ps.ByName("category"), ps.ByName("index"), ps.ByName("id"), ps.ByName("mark"))))
}

var (
	router          = httprouter.New()
	kvAddress       = flag.String("kvaddr", "", "[kvaddr] Etcd address default is empty, if empty app will not register its self.")
	enableKV        = flag.Bool("kv", false, "[kv] Enable register kv function, default false.")
	backendNameFlag = flag.String("backend", "", "[backend] Backend name.")
	watchKeyFlag    = flag.String("watch", "", "[watch] The watch path of KV.")
	showDetails     = flag.Bool("d", false, "[d] Show detail of parameter, default false.")
)

func main() {

	flag.Parse()
	printFlags()

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
	router.GET("/web_root/deep/one/two/three", Deep)
	router.GET("/web_root/deep/one/two/three_sub", Deep_Sub)
	router.GET("/web_root/restful/:category/:index/:id/:mark", Restful)

	if *enableKV {
		registerService()
	}

	log.Fatal(http.ListenAndServe("0.0.0.0:80", router))
}
func printFlags() {
	if *showDetails {
		fmt.Printf("kvaddr: %v\n", *kvAddress)
		fmt.Printf("kv: %v\n", *enableKV)
		fmt.Printf("backend: %v\n", *backendNameFlag)
		fmt.Printf("watch: %v\n", *watchKeyFlag)
	}
}

func registerService() {

	addr := *kvAddress
	backendName := *backendNameFlag
	watchKey := *watchKeyFlag

	//addr = "http://192.168.196.88:12379"
	//backendName = "backend_web_app"
	//watchKey = "/traefik/alias"

	if len(addr) != 0 {

		aliasChan := make(chan bool, 100)
		addrs := strings.Split(addr, ",")
		proxy, err := client.New(client.Config{
			Endpoints:               addrs,
			Transport:               client.DefaultTransport,
			HeaderTimeoutPerRequest: time.Second,
		})

		if err != nil {
			fmt.Printf("initial etcd failed with error: %v\n", err)
			return
		}

		currentAlias := ""
		go func() {
			keyWatchAlias := client.NewKeysAPI(proxy)
			rsp, err := keyWatchAlias.Get(context.TODO(), watchKey, &client.GetOptions{Recursive: false, Quorum: true})
			if err != nil {
				log.Println(err)
				return
			}

			currentAlias = rsp.Node.Value
			aliasChan <- true
			w := keyWatchAlias.Watcher(watchKey, &client.WatcherOptions{Recursive: false})
			for {
				rsp, err = w.Next(context.Background())
				if err != nil {
					log.Println(err)
				}

				if err == nil && currentAlias != rsp.Node.Value {
					log.Printf("Detacted node value changed from %v to %v\n", currentAlias, rsp.Node.Value)
					currentAlias = rsp.Node.Value
					aliasChan <- true
				}
			}
		}()

		go func() {

			keyKeepalive := client.NewKeysAPI(proxy)
			defaultWeight := "5"
			defaultUrl := fmt.Sprintf("http://%v:80", ip)
			defaultTTL := &client.SetOptions{TTL: time.Second * 10}

			go func() {
				for {
					time.Sleep(time.Second * 5)
					aliasChan <- true
				}
			}()

			for range aliasChan {
				if len(currentAlias) != 0 {
					keyKeepalive.Set(context.TODO(), fmt.Sprintf("%v/backends/%v/servers/%v/url", currentAlias, backendName, ip), defaultUrl, defaultTTL)
					keyKeepalive.Set(context.TODO(), fmt.Sprintf("%v/backends/%v/servers/%v/weight", currentAlias, backendName, ip), defaultWeight, defaultTTL)
				}
			}
		}()

	}
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
