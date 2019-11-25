package health

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	json "encoding/json"
	"fmt"
	"net/http"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
	"go.uber.org/zap"

	deployer_v1 "github.com/dotmesh-io/ds-deployer/apis/deployer/v1"
)

type ObjectCache interface {
	ModelDeployments() []*deployer_v1.Deployment
}

type Module interface {
	OK() bool
}

type HealthServer interface {
	RegisterModule(name string, module Module)
}

type NilHealthServer struct{}

func (s *NilHealthServer) RegisterModule(name string, module Module) {}

type Status struct {
	port       string
	httpServer *http.Server
	router     *mux.Router
	tlsConfig  *tls.Config
	logger     *zap.SugaredLogger

	username, password string

	modulesMu *sync.RWMutex
	modules   map[string]Module

	objectCache ObjectCache

	failures int
}

type Opts struct {
	Port               string
	TLSConfig          *tls.Config
	Logger             *zap.SugaredLogger
	Username, Password string
	ObjectCache        ObjectCache
}

func NewServer(opts *Opts) *Status {
	return &Status{
		port:        opts.Port,
		tlsConfig:   opts.TLSConfig,
		logger:      opts.Logger,
		username:    opts.Username,
		password:    opts.Password,
		modulesMu:   &sync.RWMutex{},
		objectCache: opts.ObjectCache,
		modules:     make(map[string]Module),
	}
}

var defaultALPN = []string{"h2", "http/1.1"}

func (s *Status) Start() error {

	n := negroni.New(negroni.NewRecovery())

	s.router = mux.NewRouter()
	s.registerRoutes(s.router)

	n.Use(negroni.HandlerFunc(s.authenticationMiddleware))
	n.Use(negroni.NewRecovery())
	n.UseHandler(s.router)

	s.logger.Infow("starting status API",
		"http_port", s.port,
	)

	if s.tlsConfig != nil {
		s.httpServer = &http.Server{
			Addr:              fmt.Sprintf(":%s", s.port),
			TLSConfig:         s.tlsConfig,
			Handler:           n,
			IdleTimeout:       time.Second * 120,
			ReadTimeout:       time.Second * 30,
			ReadHeaderTimeout: time.Second * 30,
			WriteTimeout:      time.Second * 25,
		}
		s.httpServer.TLSConfig.NextProtos = defaultALPN
		err := s.httpServer.ListenAndServeTLS("", "")
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				return nil
			}
		}
		return err
	}

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%s", s.port),
		Handler:           n,
		IdleTimeout:       time.Second * 120,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      time.Second * 25,
	}
	err := s.httpServer.ListenAndServe()
	if err != nil {
		if strings.Contains(err.Error(), "closed") {
			return nil
		}
	}
	return err
}

func (s *Status) registerRoutes(mux *mux.Router) {

	// mux.HandleFunc("/health", s.healthHandler).Methods("GET")
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/cache", s.getCacheContentsHandler).Methods("GET")
}

func (s *Status) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)

	}
	return nil
}

func (s *Status) RegisterModule(name string, module Module) {
	s.modulesMu.Lock()
	s.modules[name] = module
	s.modulesMu.Unlock()
}

//easyjson:json
type HealthResponse struct {
	Health map[string]bool `json:"health"`
}

func (s *Status) OK() bool {
	s.modulesMu.RLock()
	defer s.modulesMu.RUnlock()

	for _, v := range s.modules {
		if !v.OK() {
			return false
		}
	}

	return true
}

func (s *Status) healthHandler(w http.ResponseWriter, req *http.Request) {

	var r HealthResponse
	var down bool
	var ok bool

	s.modulesMu.RLock()

	r.Health = make(map[string]bool, len(s.modules))
	for k, v := range s.modules {
		ok = v.OK()
		r.Health[k] = ok
		if !ok {
			down = true
		}
	}

	// if we failed 5 times, fatal error
	if s.failures > 4 {
		s.logger.Fatalf("failures treshold reached")
	}

	s.modulesMu.RUnlock()

	if down {

		s.modulesMu.Lock()
		s.failures++
		s.modulesMu.Unlock()

		w.WriteHeader(http.StatusInternalServerError)
		var buf bytes.Buffer
		err := pprof.Lookup("goroutine").WriteTo(&buf, 1)

		if err != nil {
			// log.WithError(err).Error("failed to get goroutine stacktrace")
			s.logger.Errorw("failed to get goroutine stacktrace",
				"error", err,
			)
		} else {
			s.logger.Errorw("fatal error: healthcheck failed",
				"stacktrace", base64.StdEncoding.EncodeToString(buf.Bytes()),
			)
		}

	} else {

		s.modulesMu.Lock()
		s.failures = 0
		s.modulesMu.Unlock()

		w.WriteHeader(http.StatusOK)
	}

	_, err := easyjson.MarshalToWriter(&r, w)

	if err != nil {
		s.logger.Errorw("health handler: failed to marshal JSON",
			"error", err,
		)
	}
}

// authenticationMiddleware - authentication middleware placeholder
func (s *Status) authenticationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if r.URL.Path == "/health" {
		next(rw, r)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(rw, http.StatusText(404), http.StatusNotFound)
		return
	}

	if username != s.username || password != s.password {
		http.Error(rw, http.StatusText(404), http.StatusNotFound)
		return
	}

	next(rw, r)
}

func (s *Status) getCacheContentsHandler(w http.ResponseWriter, req *http.Request) {
	modelDeployments := s.objectCache.ModelDeployments()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&modelDeployments)

	if err != nil {
		s.logger.Errorw("failed to marshal cache contents",
			"error", err,
		)
		w.WriteHeader(500)
		return
	}
}
