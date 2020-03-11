package daemon

import (
	"golang_home_work/calendar/pkg/config"
	"golang_home_work/calendar/pkg/service"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run(cfg *config.Config, log *zap.Logger) error {
	var err error
	s := service.NewService(cfg.Service, log)
	if err = s.Start(); err != nil {
		return err
	}

	waitForSignal(log)

	return nil
}

func waitForSignal(log *zap.Logger) {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Info("Got signal: exiting.", zap.String("signal", s.String()))
}
