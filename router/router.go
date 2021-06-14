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
	bauth "local.com/leobrada/ztsfc_http_pep/basic_auth"
	env "local.com/leobrada/ztsfc_http_pep/env"
    //metadata "local.com/leobrada/ztsfc_http_pep/metadata"
	logwriter "local.com/leobrada/ztsfc_http_pep/logwriter"
)

type Router struct {
	tls_config *tls.Config
	frontend   *http.Server
	lw         *logwriter.LogWriter
    mode       string
    file       bool
//    md         *metadata.Cp_metadata
}

func NewRouter(_lw *logwriter.LogWriter, _mode string, _file bool) (*Router, error) {
	router := new(Router)
	router.lw = _lw
    router.mode = _mode
    router.file = _file

	router.tls_config = &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: true,
		Certificates:           nil,
		//ClientAuth:             tls.RequireAndVerifyClientCert,
		ClientAuth:				tls.VerifyClientCertIfGiven,
		ClientCAs: env.Config.CA_cert_pool_service_accepts_when_presented_by_ext,
		GetCertificate: func(cli *tls.ClientHelloInfo) (*tls.Certificate, error) {
			// load a suitable certificate that is shown to clients according the request domain/TLS SNI
            if cli.ServerName == env.Config.Sni {
                return &env.Config.X509KeyPair_presented_by_service_to_ext, nil
            }
			return nil, fmt.Errorf("Error: Could not serve a suitable certificate for %s\n", cli.ServerName)
		},
	}

	// Frontend Handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", router.ServeHTTP)
	//mux.HandleFunc("/file", router.ServeFileDownload)

	// Setting Up the Frontend Server
	router.frontend = &http.Server{
		Addr:         env.Config.Listen_addr,
		TLSConfig:    router.tls_config,
		ReadTimeout:  time.Hour * 1,
		WriteTimeout: time.Hour * 1,
		Handler:      mux,
		ErrorLog:     log.New(router.lw, "", 0),
	}

    // Create metadata
    //router.md = new(metadata.Cp_metadata)

    http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10000

	return router, nil
}

func (router *Router) SetUpSFC() bool {
	return false
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//    fmt.Printf("# of GOROUTINES: %d\n", runtime.NumGoroutine())
    // Check if user is already authenticated
    if router.mode == "direct" {
        if !bauth.User_sessions_is_valid(req) {
            if !bauth.Basic_auth(w, req) {
                return
            }
        }
    }

    if !router.file {
        fmt.Fprintf(w, "1")
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
//		service_to_add := env.Config.Service_pool[service_to_add_name]
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
//		dest, ok := env.Config.Sf_pool[sf_to_add_name]
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
//				Certificates:       []tls.Certificate{env.Config.Sf_pool[sf_to_add_name].X509KeyPair_shown_by_pep_to_sf},
//				InsecureSkipVerify: true,
//				ClientAuth:         tls.RequireAndVerifyClientCert,
//				ClientCAs:          env.Config.CA_cert_pool_pep_accepts_from_int,
//			},
//		}
//
//	} else {
//		//logr.Log_writer.Log("[ Service functions ]\n")
//		//logr.Log_writer.Log("    -\n")
//		//logr.Log_writer.Log("[ Service ]\n")
//		//logr.Log_writer.Log(fmt.Sprintf("    %s\n", service_to_add_name))
//		for _, service := range env.Config.Service_pool {
//			//		if req.TLS.ServerName == service.SNI {
//			//			proxy = httputil.NewSingleHostReverseProxy(service.Dst_url)
//			if req.TLS.ServerName == service.Sni {
//				proxy = httputil.NewSingleHostReverseProxy(service.Target_service_url)
//
//				// When the PEP is acting as a client; this defines his behavior
//				// TODO: MOVE TO A BETTER PLACE
//				proxy.Transport = &http.Transport{
//					TLSClientConfig: &tls.Config{
//						Certificates:       []tls.Certificate{env.Config.Service_pool[service_to_add_name].X509KeyPair_shown_by_pep_to_service},
//						InsecureSkipVerify: true,
//						ClientAuth:         tls.RequireAndVerifyClientCert,
//						ClientCAs:          env.Config.CA_cert_pool_pep_accepts_from_int,
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
