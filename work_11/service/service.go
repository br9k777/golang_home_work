package service

import (
	"fmt"
	"golang_home_work/work_11/config"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service struct {
	config config.ConfigService
	log    *zap.Logger
	router http.Handler
}

// NewService creates a new service
func NewService(cfg config.ConfigService, logger *zap.Logger) *Service {
	r := mux.NewRouter()
	service := &Service{cfg, logger, r}
	r.HandleFunc("/", service.ResolverHandle)
	http.Handle("/", r)

	return service
}

func (s *Service) ResolverHandle(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, `
<!DOCTYPE HTML>
<html>
	<head>
		<meta charset="utf-8">
		<title>Service</title>
	</head>
	<body>
		<center>потом тут будет сервис</center>
	</body>
</html>
`)
}

func (s *Service) Start() error {
	s.log.Info(spew.Sdump(s.config))
	s.log.Debug("Test DEBUG message")
	s.log.Info("Test INFO message")
	s.log.Warn("Test WARN message")
	s.log.Error("Test ERROR message")
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port), s.router)
}
