package router

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	logger "github.com/vs-uulm/ztsfc_http_logger"
	"github.com/vs-uulm/ztsfc_http_service/internal/app/config"
)

type Router struct {
	tlsConfig *tls.Config
	frontend  *http.Server
	sysLogger *logger.Logger
	mode      string
	file      bool
}

func NewRouter(logger *logger.Logger, mode string, file bool) (*Router, error) {
	router := new(Router)
	router.sysLogger = logger
	router.mode = mode
	router.file = file

	// Create a tls.Config struct to accept incoming connections
	router.tlsConfig = &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: false,
		Certificates:           []tls.Certificate{config.Config.X509KeyPairShownByService},
		ClientAuth:             tls.RequireAndVerifyClientCert,
		ClientCAs:              config.Config.CAcertPoolPepAcceptsFromExt,
	}

	// Frontend Handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", router.ServeHTTP)
	mux.HandleFunc("/file", router.ServeFileDownload)

	// Setting Up the Frontend Server
	router.frontend = &http.Server{
		Addr:         config.Config.Service.ListenAddr,
		TLSConfig:    router.tlsConfig,
		ReadTimeout:  time.Hour * 1,
		WriteTimeout: time.Hour * 1,
		Handler:      mux,
		ErrorLog:     log.New(router.sysLogger.GetWriter(), "", 0),
	}

	// Create metadata
	//router.md = new(metadata.Cp_metadata)

	// Packet arrival registrar
	//router.requestReception = logrus.New()
	//router.requestReception.SetLevel(logrus.InfoLevel)
	//router.requestReception.SetFormatter(&logrus.JSONFormatter{})

	// Open a file for the logger output
	//    requestReceptionLogfile, err := os.OpenFile("requests_times.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//    if err != nil {
	//        log.Fatal(err)
	//    }

	// Redirect the logger output to the file
	//router.requestReception.SetOutput(requestReceptionLogfile)

	return router, nil
}

func (router *Router) SetUpSFC() bool {
	return false
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Check if user is already authenticated
	// if router.mode == "direct" {
	// 	if !bauth.User_sessions_is_valid(req) {
	// 		if !bauth.Basic_auth(w, req) {
	// 			return
	// 		}
	// 	}
	// }

	fmt.Fprintf(w, "1")
}

func (router *Router) ServeFileDownload(w http.ResponseWriter, req *http.Request) {
	// Check if user is already authenticated
	// if router.mode == "direct" {
	// 	if !bauth.User_sessions_is_valid(req) {
	// 		if !bauth.Basic_auth(w, req) {
	// 			return
	// 		}
	// 	}
	// }

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("README.md"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, req, "./README.md")
}

func (router *Router) ListenAndServeTLS() error {
	return router.frontend.ListenAndServeTLS("", "")
}
