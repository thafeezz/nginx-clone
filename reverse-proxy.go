package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	servers      []string
	currentIndex int
	indexLock    sync.Mutex
}

func (lb *LoadBalancer) getNextServerUrl() string {
	lb.indexLock.Lock()
	defer lb.indexLock.Unlock()

	lb.currentIndex++
	lb.currentIndex = lb.currentIndex % len(lb.servers)

	return lb.servers[lb.currentIndex]
}

type ConnectionHandler struct {
	targetUrl    string
	port         int
	loadBalancer *LoadBalancer
}

func (handler *ConnectionHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	log.Printf("Handling request from client %s", req.RemoteAddr)

	serverUrl := handler.loadBalancer.getNextServerUrl()

	log.Printf("Sending request to server: %s", serverUrl)

	switch req.Method {
	case "GET":
	case "POST":
	case "PUT":
	case "DELETE":

	}
}

func main() {
	// TODO: add config file and parse?
	target := flag.String("target", "", "Url of target web server")
	portArg := flag.Int("port", 8080, "Port on which to run reverse proxy")
	servers := flag.Args()
	flag.Parse()

	numServers := len(servers)

	if numServers == 0 {
		log.Fatal("No backend servers provided")
	}

	// TODO: might not need this
	if *target == "" {
		log.Fatal("Url of target server not provided")
	}

	port := fmt.Sprintf(":%d", *portArg)

	var lb *LoadBalancer = &LoadBalancer{servers, 0, sync.Mutex{}}
	var handler *ConnectionHandler = &ConnectionHandler{targetUrl: *target, port: *portArg, loadBalancer: lb}

	log.Fatal(http.ListenAndServe(port, handler))
}
