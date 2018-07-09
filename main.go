package main


import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/cmattoon/aws-ssm/pkg/config"
	"github.com/cmattoon/aws-ssm/pkg/controller"
)


func main() {
	cfg := config.DefaultConfig()
	if err := cfg.ParseFlags(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	log.Infof("Using config: %s", cfg)
	
	stopChan := make(chan struct{}, 1)
	
	go doMetrics(cfg.MetricsListenAddress)
	go handleSigterm(stopChan)

	ctrl := controller.NewController(cfg)

	ctrl.Run(stopChan)
}

func handleSigterm(stopChan chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Info("Received SIGTERM. Terminating")
	close(stopChan)
}

func doMetrics(address string) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(address, nil))
}
