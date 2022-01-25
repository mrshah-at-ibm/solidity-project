package routes

import "net/http"

func (r *Routes) StartServer() error {
	logger := r.Logger.Named("StartServer")

	logger.Info("Starting server on port 3000")
	err := http.ListenAndServe(":3000", r.router)
	return err
}
