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

// func SetupListener() error {
// 	client, err := ethclient.Dial("http://node-1.network:32001")
// 	if err != nil {
// 		return err
// 	}

// 	contractAddress := common.HexToAddress("0xbfA009C3C51CA7718337a806754Fa0443d843913")
// 	query := ethereum.FilterQuery{
// 		Addresses: []common.Address{contractAddress},
// 	}

// 	logs := make(chan types.Log)
// 	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)

// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			return err
// 		case vLog := <-logs:
// 			fmt.Println("Received log")
// 			fmt.Println("Log: ", vLog)
// 		}
// 	}
// }

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
	// out, err := tx.MarshalBinary()
	// if err != nil {
	// 	logger.Error("Error Marshaling receipt")
	// }
	_, err = w.Write([]byte(tx.TxHash.String()))
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
	_, err = w.Write([]byte(tx.TxHash.String()))
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
	_, err = w.Write([]byte(tx.TxHash.String()))
	if err != nil {
		fmt.Println("Error sending response")
	}

	return

}

func (r *Routes) deployContractRoute(w http.ResponseWriter, req *http.Request) {
	// logger := r.Logger.Named("deployContractRoute")

	return
}

func (r *Routes) getContractRoute(w http.ResponseWriter, req *http.Request) {
	// logger := r.Logger.Named("getContractRoute")

}

// func (r *Routes) mintTokenRoute(w http.ResponseWriter, req *http.Request) {
// 	logger := r.Logger.Named("mintTokenRoute")

// 	_, err := w.Write([]byte("Placeholder mintTokenRoute"))
// 	if err != nil {
// 		logger.Error("Error sending response")
// 	}
// }

// func (r *Routes) BurnTokenRoute(w http.ResponseWriter, req *http.Request) {
// 	logger := r.Logger.Named("BurnTokenRoute")

// 	_, err := w.Write([]byte("Placeholder BurnTokenRoute"))
// 	if err != nil {
// 		logger.Error("Error sending response")
// 	}
// }

// func (r *Routes) TransferTokenRoute(w http.ResponseWriter, req *http.Request) {
// 	logger := r.Logger.Named("TransferTokenRoute")

// 	_, err := w.Write([]byte("Placeholder TransferTokenRoute"))
// 	if err != nil {
// 		logger.Error("Error sending response")
// 	}
// }

// func (r *Routes) getTokenRoute(w http.ResponseWriter, req *http.Request) {
// 	logger := r.Logger.Named("getTokenRoute")

// 	_, err := w.Write([]byte("Placeholder getTokenRoute"))
// 	if err != nil {
// 		logger.Error("Error sending response")
// 	}
// }
