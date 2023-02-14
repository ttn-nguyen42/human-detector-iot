package main

import (
	"context"
	"iot_api/apps"
	"iot_api/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(utils.GetLogLevel())
}

func main() {
	app := apps.New()
	// Runs the app in a seperate goroutine
	// The main threads manages shutdown process
	go func() {
		err := app.Start()
		if err != nil && err != http.ErrServerClosed {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Unable to start")
		}
	}()
	// Listens to user signal (Ctrl + C)
	// puts the signal info to the channel
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// Blocks until receive a signal from the channel
	<- signalChan
	logrus.Debug("Received stop signal")
	// This context includes a timeout
	// when the timeout (3s) is over, timeoutCtx.Done() channel sends a signal
	// to notify any goroutines that is listening to this channel that
	// the timeout is over
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// defer sets a function to run after the current scope has finished running
	// scope here = main()
	defer func ()  {
		logrus.Info("Cancel remaining timeout")
		timeoutCancel()
	}()
	logrus.Info("Stopping HTTP server")
	// Stops the server
	err := app.Stop(timeoutCtx)
	// Checks if context timeout is finished
	if err == context.DeadlineExceeded {
		logrus.Debug("Context timed out before all connections are closed, halted connections")
	}
	close(signalChan)
	logrus.Info("Completed graceful shutdown")
}