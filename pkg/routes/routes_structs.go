package routes

import (
	"github.com/go-chi/chi"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/executer"
	"go.uber.org/zap"
)

type Routes struct {
	router   *chi.Mux
	Logger   *zap.SugaredLogger
	Executer executer.ExecuterInterface
}
