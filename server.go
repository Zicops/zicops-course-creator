package main

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"

	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-course-creator/config"
	"github.com/zicops/zicops-course-creator/controller"
	"github.com/zicops/zicops-course-creator/global"
	"github.com/zicops/zicops-course-creator/lib/db/cassandra"
)

func main() {
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "zicops-cc.json")
	log.Infof("Starting zicops course creator service")
	ctx, cancel := context.WithCancel(context.Background())

	cassConfig1 := config.NewCassandraConfig("coursez")
	cassSession1, err := cassandra.New(cassConfig1)
	if err != nil {
		log.Errorf("Error connecting to cassandra: %s", err)
		log.Infof("zicops course creator intialization failed")
	}

	cassConfig2 := config.NewCassandraConfig("qbankz")
	cassSession2, err := cassandra.New(cassConfig2)
	if err != nil {
		log.Errorf("Error connecting to cassandra: %s", err)
		log.Infof("zicops course creator intialization failed")
	}
	log.Info("Connected to cassandra")
	global.CTX = ctx
	global.CassSession = cassSession1
	global.CassSessioQBank = cassSession2
	global.Cancel = cancel
	log.Infof("zicops course creator intialization complete")
	portFromEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portFromEnv)

	if err != nil {
		port = 8090
	}
	bootUPErrors := make(chan error, 1)
	go monitorSystem(cancel, bootUPErrors)
	controller.CCBackendController(ctx, port, bootUPErrors)
	err = <-bootUPErrors
	if err != nil {
		log.Errorf("There is an issue starting backend server for course creator: %v", err.Error())
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
	errorChannel <- fmt.Errorf("System termination signal received")
}
