package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/config"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/executer"
	"go.uber.org/zap"
)

func New(l *zap.Logger, e executer.ExecuterInterface) (*Routes, error) {

	r := Routes{
		Logger:   l.Sugar().Named("Routes"),
		Executer: e,
	}

	return &r, nil
}

func (r *Routes) SetupRoutes() {
	r.router = chi.NewRouter()
	// TODO: Look for replacing this logger with zap logger
	r.router.Use(middleware.Logger)
	r.router.Get("/", r.baseRoute)
	r.router.Route("/transaction", func(r1 chi.Router) {
		r1.Use(validateToken)
		r1.Post("/mint/{toaddress}", r.MintTransactionRoute)
		r1.Post("/token/{token}/burn", r.BurnTokenRoute)
		r1.Post("/token/{token}/transfer", r.TransferTokenRoute)

		r1.Get("/balance/{owner}", r.balanceOwnerRoute)
		// r1.Get("/token/owner/{token}", r.tokenOwnerRoute)

	})
	r.router.Route("/login", func(r1 chi.Router) {

		r1.Get("/github", r.githubLoginRoute)
		r1.Get("/github/callback", r.githubLoginCallbackRoute)
	})
}

func validateToken(next http.Handler) http.Handler {
	useauth, err := config.GetGithubClientID()
	if err != nil || useauth == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req.WithContext(req.Context()))
			return
		})
	} else {

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			token := req.Header.Get("x-auth-token")
			valid, err := config.IsTokenValid(token)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Error validating auth token"))
				return
			}

			if valid {
				next.ServeHTTP(w, req.WithContext(req.Context()))
				return
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Error validating auth token"))
				return
			}

		})
	}

}

func (r *Routes) baseRoute(w http.ResponseWriter, req *http.Request) {
	logger := r.Logger.Named("baseRoute")
	useauth, err := config.GetGithubClientID()

	fmt.Fprintf(w, "<html><body> Server running ok")
	if err != nil {
		logger.Error("Error sending response")
		return
	}

	if err != nil || useauth == "" {
		fmt.Fprintf(w, `<br/> Server running without auth`)

	} else {
		fmt.Fprintf(w, `<br/> Github auth enabled: <a href="/login/github">Generate Token</a>`)
	}

	fmt.Fprintf(w, "<body></html>")

}

func (r *Routes) MintTransactionRoute(w http.ResponseWriter, req *http.Request) {
	logger := r.Logger.Named("MintTransactionRoute")

	toaddress := chi.URLParam(req, "toaddress")
	toaddress, err := url.QueryUnescape(toaddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Incorrect to address"))
		if err != nil {
			logger.Error("Error sending response")
		}
		return
	}

	logger.Debug("Calling MintToken")
	tx, err := r.Executer.MintToken(toaddress)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Error("Error sending response")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		logger.Error("Error sending response")
	}

	return

}

func (r *Routes) balanceOwnerRoute(w http.ResponseWriter, req *http.Request) {
	logger := r.Logger.Named("balanceOwnerRoute")

	owner := chi.URLParam(req, "owner")
	owner, err := url.QueryUnescape(owner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Incorrect to address"))
		if err != nil {
			logger.Error("Error sending response")
		}
		return
	}

	logger.Debug("Calling GetBalance")
	tx, err := r.Executer.BalanceOf(owner)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Error("Error sending response")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(strconv.Itoa(tx)))
	if err != nil {
		logger.Error("Error sending response")
	}

	return

}

func (r *Routes) BurnTokenRoute(w http.ResponseWriter, req *http.Request) {
	logger := r.Logger.Named("BurnTokenRoute")

	token := chi.URLParam(req, "token")
	token, err := url.QueryUnescape(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Incorrect to address"))
		if err != nil {
			logger.Error("Error sending response")
		}
		return
	}

	logger.Debug("Calling BurnToken")
	tx, err := r.Executer.BurnToken(token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Error("Error sending response")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		logger.Error("Error sending response")
	}

	return

}

func (r *Routes) TransferTokenRoute(w http.ResponseWriter, req *http.Request) {
	logger := r.Logger.Named("TransferTokenRoute")

	type TransferStruct struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	bodybytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Error reading body"))
		if err != nil {
			logger.Error("Error sending error reading body")
		}
		return
	}

	body := TransferStruct{}

	err = json.Unmarshal(bodybytes, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Incorrect JSON"))
		if err != nil {
			logger.Error("Error sending response")
		}
		return
	}

	token := chi.URLParam(req, "token")
	token, err = url.QueryUnescape(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Incorrect to address"))
		if err != nil {
			logger.Error("Error sending response")
		}
		return
	}

	logger.Debug("Calling TransferToken")
	tx, err := r.Executer.TransferToken(body.From, body.To, token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logger.Error("Error sending response")
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		logger.Error("Error sending response")
	}

	return

}
