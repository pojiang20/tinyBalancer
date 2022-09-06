package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tinybalancer/proxy"
	"log"
	"net/http"
	"strconv"
)

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	err = config.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}

	router := mux.NewRouter()
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}
		router.Handle(l.Pattern, httpProxy)
	}
	if config.MaxAllowed > 0 {
		router.Use(maxAllowedMiddleware(config.MaxAllowed))
	}
	svr := http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(config.Port)),
		Handler: router,
	}
	config.Print()
	// listen and serve
	if config.Schema == "http" {
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if config.Schema == "https" {
		err := svr.ListenAndServeTLS(config.SSLCertificate, config.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}
