package basic_auth

import (
    "net"
	"net/http"
	"fmt"
	"crypto/rsa"
	"time"

    "gopkg.in/ldap.v2"
	"github.com/golang-jwt/jwt/v4"
    "github.com/vs-uulm/ztsfc_http_service/internal/app/config"
    logger "github.com/vs-uulm/ztsfc_http_logger"
)

func UserSessionsIsValid(req *http.Request) bool {
	jwtCookie, err := req.Cookie("ztsfc_session")
	if err != nil {
		return false
	}
	ss := jwtCookie.Value

	_, err = jwt.Parse(ss, func(token *jwt.Token) (interface{}, error) {
		return config.Config.BasicAuth.Session.JwtPubKey, nil
	})

	return err == nil
}

func BasicAuth(sysLogger *logger.Logger, w http.ResponseWriter, req *http.Request) bool {
	return performPasswdAuth(sysLogger, w, req)
}

func performPasswdAuth(sysLogger *logger.Logger, w http.ResponseWriter, req *http.Request) bool {
	var username, password string

	// TODO: Check for JW Token initially
	// Check if it is a POST request
	if req.Method == "POST" {

		if err := req.ParseForm(); err != nil {
			handleFormReponse("Parsing Error", w)
			return false
		}

		nmbr_of_postvalues := len(req.PostForm)
		if nmbr_of_postvalues != 2 {
			handleFormReponse("Wrong number of POST form values", w)
			return false
		}

		usernamel, exist := req.PostForm["username"]
		username = usernamel[0]
		if !exist {
			handleFormReponse("Username not present or wrong", w)
			return false
		}

		passwordl, exist := req.PostForm["password"]
		password = passwordl[0]
		if !exist {
			handleFormReponse("Password not present or wrong", w)
			return false
		}

		if !areUserLDAPCredentialsValid(sysLogger, username, password) {
			handleFormReponse("Authentication failed for user", w)
			return false
		}

		// Create JWT
		ss := createJWToken(config.Config.BasicAuth.Session.MySigningKey, username)

		ztsfc_cookie := http.Cookie{
			Name:   "ztsfc_session",
			Value:  ss,
			MaxAge: 1800,
			Path:   "/",
		}
		http.SetCookie(w, &ztsfc_cookie)

		// TODO: make it user configurable
		// TODO: is there a better solution for the content-length  /body length "bug"?
		req.ContentLength = 0
		http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusSeeOther) // 303
		return true

	} else {
		handleFormReponse("only post methods are accepted in this state", w)
		return false
	}
}

func handleFormReponse(msg string, w http.ResponseWriter) {
	form := `<html>
        <body>
        <center>
        <form action="/" method="post">
        <label for="fname">Username:</label>
        <input type="text" id="username" name="username"><br><br>
        <label for="lname">Password:</label>
        <input type="password" id="password" name="password"><br><br>
        <input type="submit" value="Submit">
        </form>
        </center>
        </body>
        </html>
        `

	//fmt.Println(msg)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, form)
}

func createJWToken(mySigningKey *rsa.PrivateKey, username string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		Issuer:    "ztsfc_bauth",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, _ := token.SignedString(mySigningKey)

	return ss
}

// AreUserLDAPCredentialsValid() checks a user credentials by binding to the given LDAP server
func areUserLDAPCredentialsValid(sysLogger *logger.Logger, userName, password string) bool {
        // If a user with the given name exists, obtain their full LDAP dn
        dn, ok := GetUserDNfromLDAP(sysLogger, userName)
        if !ok {
                sysLogger.Errorf("basic_auth: areUserLDAPCredentialsValid(): unable to find the user '%s'", userName)
                return false
        }

        // The user exists. Check user's password by binding to the LDAP database
        err := config.Config.BasicAuth.Ldap.LdapConn.Bind(dn, password)
        if err != nil {
                // User's password does not match
                sysLogger.Debugf("basic_auth: areUserLDAPCredentialsValid(): unable to bind with the given credentials (username='%s'): %s", userName, err.Error())
                return false
        }

        // Everything is ok
        sysLogger.Debugf("basic_auth: areUserLDAPCredentialsValid(): credentials of the user '%s' are valid", userName)
        return true
}

// GetUserDNfromLDAP() returns a user's full LDAP dn if the user's record exists in the database.
func GetUserDNfromLDAP(sysLogger *logger.Logger, userName string) (string, bool) {
        // Connect to the LDAP database with the readonly user credentials
        err := config.Config.BasicAuth.Ldap.LdapConn.Bind(config.Config.BasicAuth.Ldap.ReadonlyDN, config.Config.BasicAuth.Ldap.ReadonlyPW)
        if err != nil {
                sysLogger.Errorf("basic_auth: userNameIsInLDAP(): unable to bind to the LDAP server as the readonly user: %s", err.Error())
                return "", false
        }

        // Create a search request
        searchRequest := ldap.NewSearchRequest(
                config.Config.BasicAuth.Ldap.Base,
                ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
                fmt.Sprintf(config.Config.BasicAuth.Ldap.UserFilter, userName),
                []string{"dn"},
                nil,
        )

        fmt.Printf("LDAP REQUEST: %v\n", searchRequest)

        // Perform the search
        sr, err := config.Config.BasicAuth.Ldap.LdapConn.Search(searchRequest)
        if err != nil {
                sysLogger.Errorf("basic_auth: userNameIsInLDAP(): LDAP searching error: %s", err.Error())
                return "", false
        }

        // Nothing has been found
        if len(sr.Entries) == 0 {
                sysLogger.Debugf("basic_auth: userNameIsInLDAP(): no user '%s' in the LDAP database", userName)
                return "", false
        }

        // Too much has been found
        if len(sr.Entries) > 1 {
                sysLogger.Debugf("basic_auth: userNameIsInLDAP(): more then 1 occurence with the given filter '%s' have been found in the LDAP database", userName)
                return "", false
        }

        // Exacty what we were looking for
        sysLogger.Debugf("basic_auth: userNameIsInLDAP(): user '%s' has been found in the LDAP database", userName)
        return sr.Entries[0].DN, true
}

func FilteredByPerimter(clientReq *http.Request) bool {
    host, _, err := net.SplitHostPort(clientReq.RemoteAddr)
    if err != nil {
        return true
    }

    for _, trustedNetwork := range config.Config.BasicAuth.Perimeter.TrustedIPNetworks {
        if trustedNetwork.Contains(net.ParseIP(host)) {
            return false
        }
    }

    return true
}

// Just for LCN paper
func Perform_moodle_login(w http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("li")
	if err != nil {
		// Transform existing http request into log POST form
		//        req.Method = "POST"

		// Set cookie presenting that user is logged in
		fmt.Printf("Performing Moodle log in...\n")
		li_cookie := &http.Cookie{
			Name:   "li",
			Value:  "yes",
			MaxAge: 36000,
			Path:   "/",
		}
		//req.AddCookie(li_cookie)
		http.SetCookie(w, li_cookie)
		return true
	}

	fmt.Printf("Cookie is present, user is logged in\n")
	return false

	//    jwt_cookie, err := req.Cookie("jwt")
	//    if err != nil {
	//        return false
	//    }
	//    ss := jwt_cookie.Value
	//
	//    token, err := jwt.Parse(ss , func(token *jwt.Token) (interface{}, error) {
	//        return parseRsaPublicKeyFromPemStr("./certs/jwt_test_pub.pem"), nil
	//    })
	//
	//    if err != nil {
	//        return false
	//    }
	//
	//    aud := token.Claims.(jwt.MapClaims)["aud"]
	//    fmt.Printf("%s\n", aud)
	//    if aud == "yes" {
	//        fmt.Printf("User is already logged in.\n")
	//        return true
	//    }
	//
	//    if aud != "no" {
	//        fmt.Printf("Wrong audience value.\n")
	//        return false
	//    }
	//
	//    fmt.Printf("Performing log in for client...\n")
	//
	//    token.Claims.(jwt.MapClaims)["aud"] = "yes"
	//    mySigningKey := parseRsaPrivateKeyFromPemStr("./certs/jwt_test_priv.pem")
	//    ss, _ = token.SignedString(mySigningKey)
	//
	//    jwt_cookie.Value = ss
	//
	//    aud = token.Claims.(jwt.MapClaims)["aud"]
	//    fmt.Printf("%s\n", aud)
	//    fmt.Printf("URL: %s\n", req.URL.String())
	//    return true
}
