package main

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-course-creator/controller"
	"github.com/zicops/zicops-course-creator/global"
)

func main() {
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "zicops-cc.json")
	log.Infof("Starting zicops course creator service")
	ctx, cancel := context.WithCancel(context.Background())
	global.CTX = ctx
	global.Cancel = cancel
	log.Infof("zicops course creator intialization complete")
	portFromEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portFromEnv)

	if err != nil {
		port = 8090
	}
	gin.SetMode(gin.ReleaseMode)

	bootUPErrors := make(chan error, 1)
	go monitorSystem(cancel, bootUPErrors)
	go checkAndInitCassandraSession()
	controller.CCBackendController(ctx, port, bootUPErrors)
	err = <-bootUPErrors
	if err != nil {
		log.Errorf("there is an issue starting backend server for course creator: %v", err.Error())
		global.WaitGroupServer.Wait()
		os.Exit(1)
	}
	log.Infof("course creator server started successfully.")
}

func monitorSystem(cancel context.CancelFunc, errorChannel chan error) {
	holdSignal := make(chan os.Signal, 1)
	signal.Notify(holdSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	// if system throw any termination stuff let channel handle it and cancel
	<-holdSignal
	cancel()
	// send error to channel
	errorChannel <- fmt.Errorf("system termination signal received")
}

func checkAndInitCassandraSession() error {
	// get user session every 1 minute
	// if session is nil then create new session
	//test cassandra connection
	_, err1 := cassandra.GetCassSession("coursez")
	_, err2 := cassandra.GetCassSession("qbankz")
	_, err3 := cassandra.GetCassSession("userz")
	if err1 != nil || err2 != nil || err3 != nil {
		log.Errorf("Error connecting to cassandra: %v and %v ", err1, err2, err3)
	} else {
		log.Infof("Cassandra connection successful")
	}
	for {
		for key := range cassandra.GlobalSession {
			session, err := cassandra.GetCassSession(key)
			restart := false
			if session != nil {
				qX := session.Query("select now() from system.local", nil)
				if qX != nil {
					err = qX.Exec()
					if err != nil {
						restart = true
					}
				}
			} else {
				restart = true
			}

			if err != nil || restart {
				cassandra.GlobalSession[key] = nil
				session, err := cassandra.GetCassSession(key)
				if err == nil && session != nil {
					cassandra.GlobalSession[key] = nil
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
