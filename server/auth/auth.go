package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/neurafuse/tools-go/config"
	// infraconfig "github.com/neurafuse/tools-go/config/infrastructure"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/neurafuse/tools-go/crypto/jwt"
	"github.com/neurafuse/tools-go/errors"
	infraID "github.com/neurafuse/tools-go/infrastructures/id"
	kubeID "github.com/neurafuse/tools-go/kubernetes/client/id"
	"github.com/neurafuse/tools-go/logging"
	projectsID "github.com/neurafuse/tools-go/projects/id"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/users"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var loginUserLast string

func (f F) Check(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		logging.LogActive()
		f.Config(req)
		var headerKeyJWTToken string = "Token"
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
				if loginUserLast != usersID.F.GetActive(usersID.F{}) {
					logging.Log([]string{"", vars.EmojiCrypto, vars.EmojiSuccess}, "Login from user: "+usersID.F.GetActive(usersID.F{}), 0)
					loginUserLast = usersID.F.GetActive(usersID.F{})
				}
				endpoint(wri, req)
			}
		} else {
			var msg string = "Invalid login for user: " + usersID.F.GetActive(usersID.F{})
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
	f.setInfraIDActive(req)
	f.setClusterIDActive(req)
	if users.F.Exists(users.F{}, usersID.F.GetActive(usersID.F{})) {
		f.SetPaths()
		f.setConfigs()
		jwt.SigningKeyActive = []byte(config.Setting("get", "user", "Spec.Auth.JWT.SigningKey", ""))
	} else {
		jwt.ResetSigningKey()
	}
}

func (f F) SetPaths() {
	vars.ProjectsBasePath = users.F.GetAPIUserBasePath(users.F{})
	vars.ProjectPath = vars.ProjectsBasePath + projectsID.GetActive()
}

func (f F) setConfigs() {
	config.Setting("init", "infrastructure", "", "")
}

func (f F) setUsersIDActive(req *http.Request) {
	if len(req.Header["From"]) >= 1 {
		usersID.F.SetActive(usersID.F{}, req.Header["From"][0])
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received request without valid user origin!", true, false, true)
	}
}

func (f F) setProjectIDActive(req *http.Request) {
	var projectID string = f.CookieReader(req.Header["Cookie"], "projectID")
	if projectID != "" {
		projectsID.SetActive(projectID)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received request without valid projectID!", true, false, true)
	}
}

func (f F) setInfraIDActive(req *http.Request) {
	var id string = f.CookieReader(req.Header["Cookie"], "infraID")
	if id != "" {
		infraID.F.SetActive(infraID.F{}, id)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received request without valid infraID!", true, false, true)
	}
}

func (f F) setClusterIDActive(req *http.Request) {
	var id string = f.CookieReader(req.Header["Cookie"], "kubeID")
	if id != "" {
		kubeID.F.SetActive(kubeID.F{}, id)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received request without valid kubeID!", true, false, true)
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
