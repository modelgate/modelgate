package server

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	StatusStarted = 1
	StatusStoped  = 2
)

type Server struct {
	profile *Profile
	server  *http.Server
	status  atomic.Uint32
	mu      sync.RWMutex
}

func NewH2CServer(profile *Profile) (s *Server, err error) {
	handler := registerHandlers(profile.Container)
	return newServer(profile, &http.Server{
		Addr:    profile.Addr,
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}), nil
}

func NewServer(profile *Profile) (s *Server, err error) {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           2 * time.Hour,
	}))

	registerRouters(profile.Container, engine)

	return newServer(profile, &http.Server{
		Addr:    profile.Addr,
		Handler: engine,
	}), nil
}

func newServer(profile *Profile, server *http.Server) *Server {
	return &Server{
		profile: profile,
		server:  server,
	}
}

func (s *Server) Start() error {
	if !s.tryStart() {
		log.Infof("%s server stoped, skip start", s.profile.Name)
		return nil
	}

	log.Infof("%s server started at %s", s.profile.Name, s.profile.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	if !s.tryStop() {
		log.Infof("%s server not started, skip shutdown", s.profile.Name)
		return
	}

	log.Infof("%s server shutdown", s.profile.Name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("failed to shutdown server: %s, error: %v\n", s.profile.Name, err)
	}
}

// 执行过 stop 就无法再次启动
// 使用 goroutine 启动多个服务是，如果一个服务启动失败，其他所有服务都会标记为 stoped。并且确保 stoped 的服务，无法被启动
func (s *Server) tryStart() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Load() == StatusStoped {
		return false
	}
	s.status.Store(StatusStarted)
	return true
}

// 如果当前状态是started, 则需要执行后续 shutdown 操作
// 无论当前是否是 started,都需要将状态标记为 stoped，防止其他 goroutine 中并发启动服务
func (s *Server) tryStop() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Load() == StatusStarted {
		s.status.Store(StatusStoped)
		return true
	}
	s.status.Store(StatusStoped)
	return false
}
