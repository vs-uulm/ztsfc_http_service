package router

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//    "runtime"
	//	"net/http/httputil"
	"time"

	logger "github.com/vs-uulm/ztsfc_http_logger"
	bauth "github.com/vs-uulm/ztsfc_http_service/internal/app/basic_auth"
	"github.com/vs-uulm/ztsfc_http_service/internal/app/config"
	//metadata "local.com/leobrada/ztsfc_http_pep/metadata"
	// PACKET ARRIVAL
	//    logrus "github.com/sirupsen/logrus"
)

type Router struct {
	tlsConfig *tls.Config
	frontend  *http.Server
	sysLogger *logger.Logger
//	mode      string
//	file      bool
	//    md         *metadata.Cp_metadata
	// PACKET ARRIVAL
	//requestReception *logrus.Logger
}

//func NewRouter(logger *logger.Logger, mode string, file bool) (*Router, error) {
func NewRouter(logger *logger.Logger, mode string, file bool) (*Router, error) {
	router := new(Router)
	router.sysLogger = logger
//	router.mode = mode
//	router.file = file

	// Create a tls.Config struct to accept incoming connections
	router.tlsConfig = &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: false,
		Certificates:           []tls.Certificate{config.Config.Service.X509KeyPairShownByService},
		ClientAuth:             tls.RequireAndVerifyClientCert,
		ClientCAs:              config.Config.Service.CAcertPoolServiceAcceptsFromExt,
	}

	// Frontend Handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", router.ServeHTTP)
	//mux.HandleFunc("/file", router.ServeFileDownload)

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
	if config.Config.Service.Mode == "direct" {
        if config.Config.BasicAuth.Perimeter.ApplyPerimeter && bauth.FilteredByPerimter(req) {
            router.sysLogger.Infof("Request from '%s' has been dropped due to a perimeter filter rule", req.RemoteAddr)
            return
        }

		if !bauth.UserSessionsIsValid(req) {
			if !bauth.BasicAuth(router.sysLogger, w, req) {
				return
			}
		}
	}

	if !config.Config.Service.File {
		fmt.Fprintf(w, webpage)
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("bigfile"))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, req, "./bigfile")
	}
}

//func (router *Router) ServeFileDownload(w http.ResponseWriter, req *http.Request) {
//    // Check if user is already authenticated
//    if !bauth.User_sessions_is_valid(req) {
//        if !bauth.Basic_auth(w, req) {
//            return
//        }
//    }
//
//    w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("README.md"))
//    w.Header().Set("Content-Type", "application/octet-stream")
//    http.ServeFile(w, req, "./README.md")
//}

//func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	// Log all http requests incl. TLS information
//	logwriter.Log_writer.LogHTTPRequest(req)
//
//	// Check if the user is authenticated; if not authenticate him/her; if that fails return an error
//	// TODO: return error to client?
//    // Check if user is already authenticated
//    if !bauth.User_sessions_is_valid(req) {
//        if !bauth.Basic_auth(w, req) {
//            return
//        }
//    }
//
//    // Authorization
//    // fmt.Printf("An dieser Stelle könnte ihre Authorisierung stehen.\n")
//
//    // Authentication Perform Login In
//    //if !bauth.Perform_moodle_login(req) {
//    //    fmt.Printf("User not logged in.\n")
//    //}
//    bauth.Perform_moodle_login(w, req)
//
//	// If user could be authenticated, create ReverseProxy variable for the connection to serve
//	var proxy *httputil.ReverseProxy
//
//	// ===== GARBAGE STARTING FROM HERE =====
//
//	// HE COMES THE LOGIC IN THIS FUNCTION
//	need_to_go_through_sf := router.SetUpSFC()
//
//	// Forward packets through the SF "Logger"
//	need_to_go_through_logger := true
//
//	// need_to_go_through_sf = false
//
//	sf_to_add_name := "dummy"
//	service_to_add_name := "nginx"
//
//	if need_to_go_through_sf {
//		/*
//		   Here comes a Magic:
//		   Definition a set of Sfs to go through
//		   ...
//
//		   Adding SF information to the HTTP header
//		   ...
//		*/
//
//		//logr.Log_writer.Log("[ Service functions ]\n")
//		//logr.Log_writer.Log(fmt.Sprintf("    - %s\n", sf_to_add_name))
//		//logr.Log_writer.Log("[ Service ]\n")
//		//logr.Log_writer.Log(fmt.Sprintf("    %s\n", service_to_add_name))
//
//		// Temporary Solution
//		service_to_add := config.Config.Service_pool[service_to_add_name]
//		/*
//		   req.Header.Add("service", service_to_add.Dst_url.String())
//		*/
//		// TODO CRUCIAL: Delete existing SFP headers for security reasons.
//		sfp, ok := req.Header["Sfp"]
//		if ok {
//			req.Header.Del("Sfp")
//		}
//		sfp = append(sfp, service_to_add.Target_service_addr)
//		req.Header["Sfp"] = sfp
//
//		// Set the SF "Logger" verbosity level
//		if need_to_go_through_logger {
//			LoggerHeaderName := "Sfloggerlevel"
//			_, ok := req.Header[LoggerHeaderName]
//			if ok {
//				req.Header.Del(LoggerHeaderName)
//			}
//
//			req.Header[LoggerHeaderName] = []string{fmt.Sprintf("%d",
//				// logwriter.SFLOGGER_REGISTER_PACKETS_ONLY |
//				logwriter.SFLOGGER_PRINT_GENERAL_INFO|
//					logwriter.SFLOGGER_PRINT_HEADER_FIELDS|
//					// logwriter.SFLOGGER_PRINT_BODY|
//					// logwriter.SFLOGGER_PRINT_FORMS|
//					// logwriter.SFLOGGER_PRINT_FORMS_FILE_CONTENT|
//					// logwriter.SFLOGGER_PRINT_TRAILERS|
//					//logwriter.SFLOGGER_PRINT_TLS_MAIN_INFO|
//					//logwriter.SFLOGGER_PRINT_TLS_CERTIFICATES|
//					// logwriter.SFLOGGER_PRINT_TLS_PUBLIC_KEY |
//					// logwriter.SFLOGGER_PRINT_TLS_CERT_SIGNATURE |
//					// logwriter.SFLOGGER_PRINT_RAW |
//					// logwriter.SFLOGGER_PRINT_REDIRECTED_RESPONSE|
//					// logwriter.SFLOGGER_PRINT_EMPTY_FIELDS |
//					0)}
//		}
//
//		dest, ok := config.Config.Sf_pool[sf_to_add_name]
//		if !ok {
//			w.WriteHeader(503)
//			return
//		}
//		proxy = httputil.NewSingleHostReverseProxy(dest.Target_sf_url)
//
//		proxy.ErrorLog = log.New(router.lw, "", 0)
//
//		// When the PEP is acting as a client; this defines his behavior
//		proxy.Transport = &http.Transport{
//			TLSClientConfig: &tls.Config{
//				// TODO: Replace it by loading the cert for the first SF in the chain
//				Certificates:       []tls.Certificate{config.Config.Sf_pool[sf_to_add_name].X509KeyPair_shown_by_pep_to_sf},
//				InsecureSkipVerify: true,
//				ClientAuth:         tls.RequireAndVerifyClientCert,
//				ClientCAs:          config.Config.CA_cert_pool_pep_accepts_from_int,
//			},
//		}
//
//	} else {
//		//logr.Log_writer.Log("[ Service functions ]\n")
//		//logr.Log_writer.Log("    -\n")
//		//logr.Log_writer.Log("[ Service ]\n")
//		//logr.Log_writer.Log(fmt.Sprintf("    %s\n", service_to_add_name))
//		for _, service := range config.Config.Service_pool {
//			//		if req.TLS.ServerName == service.SNI {
//			//			proxy = httputil.NewSingleHostReverseProxy(service.Dst_url)
//			if req.TLS.ServerName == service.Sni {
//				proxy = httputil.NewSingleHostReverseProxy(service.Target_service_url)
//
//				// When the PEP is acting as a client; this defines his behavior
//				// TODO: MOVE TO A BETTER PLACE
//				proxy.Transport = &http.Transport{
//					TLSClientConfig: &tls.Config{
//						Certificates:       []tls.Certificate{config.Config.Service_pool[service_to_add_name].X509KeyPair_shown_by_pep_to_service},
//						InsecureSkipVerify: true,
//						ClientAuth:         tls.RequireAndVerifyClientCert,
//						ClientCAs:          config.Config.CA_cert_pool_pep_accepts_from_int,
//					},
//				}
//			} else {
//				w.WriteHeader(503)
//				return
//			}
//		}
//	}
//
//	// ======= END GARBAGE =======
//
//	proxy.ServeHTTP(w, req)
//}

func (router *Router) ListenAndServeTLS() error {
	return router.frontend.ListenAndServeTLS("", "")
}
