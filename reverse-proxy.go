package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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

func connectionHandler(targetUrl string, port int, lb *LoadBalancer, conn net.Conn) {
	defer conn.Close()

	serverUrl := lb.getNextServerUrl()

	buf := make([]byte, 1024)

	bytesRead, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Failed to read bytes") // TODO: fatal or just log?
	}

	log.Println(bytesRead)

	log.Printf("Sending request to server: %s", serverUrl)

}

func main() {
	// TODO: add config file and parse?
	target := flag.String("target", "", "Url of target web server")
	portArg := flag.Int("port", 8080, "Port on which to run reverse proxy")
	flag.Parse()

	servers := flag.Args()
	numServers := len(servers)

	// TODO: might not need this
	if *target == "" {
		log.Fatal("Url of target server not provided")
	}

	if numServers == 0 {
		log.Fatal("No backend servers provided")
	}

	port := fmt.Sprintf(":%d", *portArg)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s", port)
	}

	lb := &LoadBalancer{servers, 0, sync.Mutex{}}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection")
		}

		go connectionHandler(*target, *portArg, lb, conn)
	}
}
