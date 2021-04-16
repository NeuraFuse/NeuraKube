package auth

import (
	"fmt"
	"net/http"
	"strings"

	"../../../tools-go/config"
	infraconfig "../../../tools-go/config/infrastructure"
	"../../../tools-go/crypto/jwt"
	"../../../tools-go/errors"
	"../../../tools-go/logging"
	"../../../tools-go/projects"
	"../../../tools-go/runtime"
	"../../../tools-go/users"
	"../../../tools-go/vars"
	gconfig "../../infrastructure/providers/gcloud/config"
	"../routines"
	jwtgo "github.com/dgrijalva/jwt-go"
)

type F struct{}

var loginUserLast string = ""

func (f F) Check(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		logging.LogActive()
		f.Config(req)
		headerKeyJWTToken := "Token"
		if req.Header[headerKeyJWTToken] != nil {
			token, err := jwtgo.Parse(req.Header[headerKeyJWTToken][0], func(token *jwtgo.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(runtime.F.GetCallerInfo(runtime.F{}, true) + ": An error occured during JWT token parsing!")
				}
				return jwt.SigningKeyActive, nil
			})
			if err != nil {
				fmt.Fprintf(wri, err.Error())
			}
			if token.Valid {
				if loginUserLast != users.GetIDActive() {
					logging.Log([]string{"", vars.EmojiCrypto, vars.EmojiSuccess}, "Login from user: "+users.GetIDActive(), 0)
					loginUserLast = users.GetIDActive()
				}
				endpoint(wri, req)
			}
		} else {
			msg := "Invalid login for user: " + users.GetIDActive()
			logging.Log([]string{"", vars.EmojiError, vars.EmojiCrypto}, msg, 0)
			fmt.Fprintf(wri, msg)
		}
	})
}

func (f F) NoAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		endpoint(wri, req)
	})
}

func (f F) Config(req *http.Request) {
	f.setUsersIDActive(req)
	f.setProjectIDActive(req)
	if users.Exists(users.GetIDActive()) {
		f.SetPaths()
		f.setConfigs()
		routines.F.Router(routines.F{})
		jwt.SigningKeyActive = []byte(config.Setting("get", "user", "Spec.Auth.JWT.SigningKey", ""))
	}
}

func (f F) SetPaths() {
	vars.ProjectsBasePath = users.BasePath + "/" + users.GetIDActive()
	config.Setting("get", "cli", "Spec.Projects.ActiveID", "") = projects.IDActive
	vars.ProjectPath = vars.ProjectsBasePath + config.Setting("get", "cli", "Spec.Projects.ActiveID", "")
}

func (f F) setConfigs() {
	config.Setting("init", "infrastructure", "", "")
	if config.ValidSettings("infrastructure", "kube", false) {
		infraconfig.F.SetKubeConfig(infraconfig.F{})
	}
	if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, false) {
		gconfig.F.SetConfigs(gconfig.F{})
	}
}

func (f F) setUsersIDActive(req *http.Request) {
	if len(req.Header["From"]) >= 1 {
		users.SetIDActive(req.Header["From"][0])
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received request without valid user origin!", true, false, true)
	}
}

func (f F) setProjectIDActive(req *http.Request) {
	projectID := f.CookieReader(req.Header["Cookie"], "projectID")
	if projectID != "" {
		projects.IDActive = projectID
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received request without valid project selection!", true, false, true)
	}
}

func (f F) CookieReader(cookies []string, key string) string {
	if len(cookies) != 0 {
		cookie := cookies[0]
		var entries []string
		if strings.Count(cookie, ";") > 1 {
			entries = strings.Split(cookie, "; ")
		} else {
			entries = []string{strings.Trim(cookie, ";")}
		}
		if len(entries) != 0 {
			m := make(map[string]string)
			for _, e := range entries {
				parts := strings.Split(e, "=")
				m[parts[0]] = parts[1]
			}
			if val, ok := m[key]; ok {
				return val
			} else {
				return ""
			}
		} else {
			return ""
		}
	} else {
		return ""
	}
}
