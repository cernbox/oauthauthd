package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/diocas/oauthauthd/pkg"

	"go.uber.org/zap"
)

// func BasicAuthOnly(logger *zap.Logger, userBackend pkg.UserBackend, sleepPause int) http.Handler {
// 	validBasicAuthsCounter := prometheus.NewCounter(prometheus.CounterOpts{
// 		Name: "valid_auths_basic",
// 		Help: "Number of valid authentications using basic authentication.",
// 	})
// 	invalidBasicAuthsCounter := prometheus.NewCounter(prometheus.CounterOpts{
// 		Name: "invalid_auths_basic",
// 		Help: "Number of valid authentications using basic authentication.",
// 	})

// 	prometheus.Register(validBasicAuthsCounter)
// 	prometheus.Register(invalidBasicAuthsCounter)

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// // Create return string
// 		// var request []string
// 		// // Add the request string
// 		// url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
// 		// request = append(request, url)
// 		// // Add the host
// 		// request = append(request, fmt.Sprintf("Host: %v", r.Host))
// 		// // Loop through headers
// 		// for name, headers := range r.Header {
// 		// 	name = strings.ToLower(name)
// 		// 	for _, h := range headers {
// 		// 		request = append(request, fmt.Sprintf("%v: %v", name, h))
// 		// 	}
// 		// }

// 		// // If this is a POST, add post data
// 		// if r.Method == "POST" {
// 		// 	r.ParseForm()
// 		// 	request = append(request, "\n")
// 		// 	request = append(request, r.Form.Encode())
// 		// }

// 		// logger.Info("REQUEST ", zap.String(r.Method, strings.Join(request, " /// ")))

// 		w.Header().Set("WWW-Authenticate", "Bearer realm='ownCloud'")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		// w.WriteHeader(http.StatusOK)
// 		return

// 		// path := r.Header.Get("auth-path")
// 		// token := r.Header.Get("auth-token")

// 		// if path == "" || token == "" {
// 		// 	invalidBasicAuthsCounter.Inc()
// 		// 	logger.Info("MISSING HEADERS")
// 		// 	w.WriteHeader(http.StatusUnauthorized)
// 		// 	return
// 		// }

// 		// path = filepath.Clean(path)

// 		// err := userBackend.Authenticate(r.Context(), path, token)
// 		// if err != nil {
// 		// 	invalidBasicAuthsCounter.Inc()
// 		// 	logger.Info("WRONG PATH OR TOKEN")
// 		// 	w.WriteHeader(http.StatusUnauthorized)
// 		// 	return
// 		// }

// 		// path_components := strings.Split(path, "/")
// 		// var user string

// 		// // Assuming EOS username is always a subdirectory of its 1st letter
// 		// // Otherwise we need to remove all possible EOS base paths
// 		// for i, elem := range path_components {
// 		// 	if len(elem) == 1 {
// 		// 		if i + 1 < len(path_components) {
// 		// 			user = path_components[i + 1]
// 		// 		}
// 		// 		break
// 		// 	}
// 		// }

// 		// validBasicAuthsCounter.Inc()
// 		// logger.Info("AUTHENTICATION SUCCEEDED", zap.String("PATH", path))
// 		// w.Header().Set("user", user)
// 		// w.WriteHeader(http.StatusOK)
// 	})
// }

func Status(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"installed\":true,\"maintenance\":false,\"needsDbUpgrade\":false,\"version\":\"10.0.10.4\",\"versionstring\":\"10.0.10\",\"edition\":\"Community\",\"productname\":\"ownCloud\"}"))
		w.WriteHeader(http.StatusOK)
	})
}

func BasicAuthOnly(logger *zap.Logger, userBackend pkg.UserBackend, sleepPause int) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqToken := r.Header.Get("Authorization")

		if reqToken != "" {
			splitToken := strings.Split(reqToken, "Bearer")
			reqToken = splitToken[1]
		}

		if reqToken == "" {
			logger.Info("NO TOKEN PROVIDED")
			time.Sleep(time.Second * time.Duration(sleepPause))
			w.Header().Set("WWW-Authenticate", "Bearer realm='ownCloud'")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := userBackend.Authenticate(r.Context(), reqToken)
		if err != nil {
			logger.Error("AUTHENTICATION FAILED", zap.Error(err), zap.String("TOKEN", reqToken))
			w.Header().Set("WWW-Authenticate", "Bearer realm='ownCloud'")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Info("AUTHENTICATION SUCCEDED", zap.String("USERNAME", user))
		w.Header().Set("user", user)
		w.WriteHeader(http.StatusOK)
	})
}
