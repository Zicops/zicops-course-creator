package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/global"
	graceful "gopkg.in/tylerb/graceful.v1" // see: https://github.com/tylerb/graceful
)

// CCBackendController ....
func CCBackendController(ctx context.Context, port int, errorChannel chan error) {
	log.Infof("Initializing router and endpoints.")
	ccRouter, err := CCRouter()
	if err != nil {
		errorChannel <- err
		return
	}
	httpAddress := fmt.Sprintf(":%d", port)
	global.WaitGroupServer.Add(1)
	go serverHTTPRoutes(ctx, httpAddress, ccRouter, errorChannel)
}

func serverHTTPRoutes(ctx context.Context, httpAddress string, handler http.Handler, errorChannel <-chan error) {
	defer global.WaitGroupServer.Done()
	// init graceful server
	serverGrace := &graceful.Server{
		Timeout: 10 * time.Second,
		//BeforeShutdown:    beforeShutDown,
		ShutdownInitiated: shutDownBackend,
		Server: &http.Server{
			Addr:    httpAddress,
			Handler: handler,
			WriteTimeout: 50 * time.Second,
			ReadTimeout:  50 * time.Second,
		},
	}
	stopChannel := serverGrace.StopChan()
	err := serverGrace.ListenAndServe()
	if err != nil {
		log.Fatalf("CCController: Failed to start server : %s", err.Error())
	}
	log.Infof("Backend is serving the routes.")
	for {
		// wait for the server to stop or be canceled
		select {
		case <-stopChannel:
			log.Infof("CCController: Server shutdown at %s", time.Now())
			return
		case <-ctx.Done():
			log.Infof("CCController: context done is called %s", time.Now())
			serverGrace.Stop(time.Second * 2)
		}
	}
}

func shutDownBackend() {
	log.Infof("CCController: Shutting down server at %s", time.Now())
}
