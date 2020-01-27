package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/grubastik/kubernetes-admission-control/internal/addlabels"
	"github.com/grubastik/kubernetes-admission-control/internal/root"
	"github.com/julienschmidt/httprouter"
)

func Init(tls *tls.Config) *http.Server {
	//configure router
	router := httprouter.New()
	router.GET("/", root.Index)
	router.POST("/append-label", addlabels.AppendLabelsHandler)

	//configure server
	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         "0.0.0.0:8443",
		TLSConfig:    tls,
		Handler:      router,
	}

	return s
}
