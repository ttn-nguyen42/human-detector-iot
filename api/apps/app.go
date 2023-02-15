package apps

import (
	"context"
	"fmt"
	"iot_api/database"
	"iot_api/routes"
	"iot_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App interface {
	Start() error
	Stop(timeoutCtx context.Context) error
}

type app struct {
	Server *http.Server
}

func New() App {
	router := gin.New()
	routes.Create(router)

	port := utils.GetPort()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	return &app{
		Server: server,
	}
}

func (a *app) Start() error {
	return a.Server.ListenAndServe()
}

func (a *app) Stop(timeoutCtx context.Context) error {
	database.Close()
	err := a.Server.Shutdown(timeoutCtx)
	return err
}
