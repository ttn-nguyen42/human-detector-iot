package apps

import (
	"context"
	"fmt"
	"iot_api/database"
	"iot_api/routes"
	"iot_api/utils"
	"iot_api/workers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App interface {
	Start() error
	Stop() error
}

type app struct {
	Server *http.Server
	Dispatcher *workers.Dispatcher
}

func New() app {
	router := gin.New()
	dispatcher := workers.NewDispatcher(10)

	// Setup API routes
	routes.Create(router)

	port := utils.GetPort()

	server := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: router,
	}

	return app{
		Server: server,
		Dispatcher: dispatcher,
	}
}

func (a *app) Start() error {
	a.Dispatcher.Breath()
	return a.Server.ListenAndServe()
}

func (a *app) Stop(timeoutCtx context.Context) error {
	a.Dispatcher.Kill()
	database.Close()
	err := a.Server.Shutdown(timeoutCtx)
	return err
}
