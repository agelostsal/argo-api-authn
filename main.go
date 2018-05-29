package main

import (
	"net/http"

	"github.com/gorilla/handlers"

	"fmt"

	"flag"

	"crypto/tls"

	"strconv"

	"github.com/ARGOeu/argo-api-authn/auth"
	"github.com/ARGOeu/argo-api-authn/config"
	"github.com/ARGOeu/argo-api-authn/routing"
	"github.com/ARGOeu/argo-api-authn/stores"
	log "github.com/Sirupsen/logrus"
)

func main() {

	// Retrieve configuration file location through cmd argument
	var cfgPath = flag.String("config", "/etc/argo-api-authN/argo-api-authn-config.json", "Path for the required configuration file.")
	flag.Parse()

	// initialize the config
	var cfg = &config.Config{}
	if err := cfg.ConfigSetUp(*cfgPath); err != nil {
		log.Error(err.Error())
		panic(err.Error())
	}

	//configure datastore
	store := &stores.MongoStore{
		Server:   cfg.MongoHost,
		Database: cfg.MongoDB,
	}
	store.SetUp()

	// configure the TLS config for the server
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS10,
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  auth.Load_CAs(cfg.CertificateAuthorities),
	}

	api := routing.NewRouting(routing.ApiRoutes, store, cfg)

	log.Info(fmt.Sprintf("%+v", cfg))
	xReqWithConType := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-AuthModel"})
	allowVerbs := handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "PUT", "DELETE", "HEAD"})

	server := &http.Server{
		Addr:      ":" + strconv.Itoa(cfg.ServicePort),
		Handler:   handlers.CORS(xReqWithConType, allowVerbs)(api.Router),
		TLSConfig: tlsConfig,
	}

	//Start the server
	err := server.ListenAndServeTLS(cfg.Certificate, cfg.CertificateKey)
	if err != nil {
		log.Fatal("API", "\t", "ListenAndServe:", err)
	}
}
