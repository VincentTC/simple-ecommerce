package webserver

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/VincentTC/simple-ecommerce/util/router"

	"golang.org/x/net/netutil"
)

// Options for the web Handler.
type Options struct {
	ListenAddress   string        `envconfig:"LISTEN_ADDRESS"`
	MaxConnections  int           `envconfig:"MAX_CONNECTION"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT"`
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT"`
	AllowedCors     string        `envconfig:"ALLOWED_CORS"`
	AllowedHeaders  string        `envconfig:"ALLOWED_HEADERS"`
	AppIdMap        map[string]int
}

// Server serves various HTTP endpoints of the application server
type Server interface {
	Run() chan error
	RunGraceful() error
	Router() *router.Router
	Stop() error
}

type server struct {
	router  *router.Router
	server  *http.Server
	options *Options
}

// New initializes a new web Handler.
func New(options *Options) Server {
	return NewWithEngine(router.HTTPRouter, options)
}

// NewWithEngine ...
func NewWithEngine(engine router.EngineType, options *Options) Server {
	arrCors := strings.Split(options.AllowedCors, "|")
	router := router.NewWithEngine(engine, arrCors, options.AllowedHeaders, options.AppIdMap)
	sv := server{}
	sv.router = router
	sv.options = options

	router.GET("/ping", sv.ping)
	router.GET("/healthz", sv.healthz)

	return &sv
}

func (s *server) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (s *server) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Router of web Handler.
func (s *server) Router() *router.Router {
	return s.router
}

// RunGraceful run the webserver with blocking
func (s *server) RunGraceful() error {
	// step 1. Run the server async
	ch := s.Run()

	// step 2. Listen to the signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// step 3. if there is a signal, then graceful stop the server
	select {
	case err := <-ch:
		return err
	case <-sigChan:
	}

	return s.Stop()
}

// Run serves the HTTP endpoints.
func (s *server) Run() chan error {
	var ch = make(chan error)
	go s.run(ch)
	return ch
}

func (s *server) run(ch chan error) {
	listener, err := net.Listen("tcp", s.options.ListenAddress)
	if err != nil {
		ch <- err
		return
	}

	if s.options.MaxConnections > 0 {
		listener = netutil.LimitListener(listener, s.options.MaxConnections)
	}

	s.server = &http.Server{
		Handler:      s.router,
		ReadTimeout:  s.options.ReadTimeout,
		WriteTimeout: s.options.WriteTimeout,
	}
	ch <- s.server.Serve(listener)
}

// Stop terminate the server gracefully
func (s *server) Stop() error {
	if s.server == nil {
		return nil
	}

	timeout := s.options.GracefulTimeout
	if timeout <= 0 {
		timeout = time.Second * 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
