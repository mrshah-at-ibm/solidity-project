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
	"github.com/mrshah-at-ibm/kaleido-project/pkg/executer"
	"go.uber.org/zap"
)

func New(l *zap.Logger, e executer.ExecuterInterface) (*Routes, error) {

	r := Routes{
		Logger: l.Sugar().Named("Routes"),
	}

	// e, err := executer.NewExecuter(l)
	// if err != nil {
	// 	r.Logger.Errorw("Error creating executer", err)
	// 	return nil, err
	// }
	r.Executer = e

	// e.DeployContract()

	// opts := bind.CallOpts{
	// 	// Pending: "",
	// 	// From: "",
	// 	// BlockNumber: "",
	// 	Context: context.TODO(),
	// }
	// symbol, err := mrst.Symbol(&opts)
	// fmt.Println("Symbol: ':", symbol, "'")

	return &r, nil
}

func (r *Routes) SetupRoutes() {
	r.router = chi.NewRouter()
	// TODO: Look for replacing this logger with zap logger
	r.router.Use(middleware.Logger)
	r.router.Get("/", r.baseRoute)
	r.router.Route("/transaction", func(r1 chi.Router) {
		r1.Post("/mint/{toaddress}", r.MintTransactionRoute)
		r1.Post("/token/{token}/burn", r.BurnTokenRoute)
		r1.Post("/token/{token}/transfer", r.TransferTokenRoute)

		r1.Get("/balance/{owner}", r.balanceOwnerRoute)
		// r1.Get("/token/owner/{token}", r.tokenOwnerRoute)

	})
}

func (r *Routes) baseRoute(w http.ResponseWriter, req *http.Request) {
	logger := r.Logger.Named("baseRoute")

	_, err := w.Write([]byte("Server running ok"))
	if err != nil {
		logger.Error("Error sending response")
	}
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
		fmt.Println("Error sending response")
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
		fmt.Println("Error sending response")
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
		fmt.Println("Error sending response")
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
		fmt.Println("Error sending response")
	}

	return

}
