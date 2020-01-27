package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/grubastik/kubernetes-admission-control/app/server"
)

func main() {
	log.Println("Starting app")
	//configure waitgroup
	var wg sync.WaitGroup
	wg.Add(1)

	//configure signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//prepare certificate
	cert, err := tls.LoadX509KeyPair("certs/webhook-append-label.pem", "certs/webhook-append-label.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	tls := tls.Config{Certificates: []tls.Certificate{cert}}

	s := server.Init(&tls)

	go func() {
		// mark gorouting as finished
		defer wg.Done()
		//start to serve requests
		log.Println("Start to serve queries")
		err := s.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	//got INT or TERM signal(starting graceful shutdown)
	<-sigs
	//close server
	s.Close()
	//wait for server to quit
	wg.Wait()
}
