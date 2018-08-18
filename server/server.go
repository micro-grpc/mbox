package server

import (
	tlspkg "crypto/tls"
	"crypto/x509"
	"fmt"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcvalidtor "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	cr "github.com/micro-grpc/mbox/consul"
	"github.com/micro-grpc/mbox/lib"
	"github.com/micro-grpc/mbox/services/authorize"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

// Server structure  gRPC server
type Server struct {
	Name              string
	Namespace         string
	Domain            string
	Description       string
	Version           string
	Conn              *grpc.Server
	Addr              string
	port              int
	ip                string
	HealthChecks      []string
	serverName        string
	ServiceID         string
	ConsulAddress     string
	consulTtl         int
	consulInterval    time.Duration
	Critical          uint
	verbose           int
	debug             bool
	prom              bool
	openTracing       bool
	isAuth            bool
	validation        bool
	recovery          bool
	isLogrus          bool
	isConsul          bool
	tls               bool
	caFile            string
	certFile          string
	keyFile           string
	limiter           *rate.Limiter
	maxRetries        uint
	promHistograms    bool
	LogrusEntry       *log.Entry
	extract           bool
	secret            string
	security          bool
	schemas           []string
	AuthFunc          *grpcauth.AuthFunc
	Authorize         *authorize.Authorize
	middlewareStreams []grpc.StreamServerInterceptor
	middlewareUnarys  []grpc.UnaryServerInterceptor
}

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type ServerOption func(*Server)

// NewServer creates a gRPC server which has no service registered and has not
// started to accept requests yet.
func NewServer(name string, namespace string, ip string, port int, options ...ServerOption) (server *Server) {
	server = &Server{
		Version:        "latest",
		HealthChecks:   nil,
		Critical:       0,
		validation:     true,
		recovery:       true,
		verbose:        0,
		debug:          false,
		prom:           false,
		openTracing:    false,
		isAuth:         false,
		isLogrus:       false,
		isConsul:       false,
		tls:            false,
		serverName:     "",
		Name:           name,
		Namespace:      namespace,
		Domain:         "",
		Description:    "",
		ConsulAddress:  "",
		caFile:         "",
		certFile:       "",
		keyFile:        "",
		limiter:        rate.NewLimiter(rate.Limit(1000), 10),
		maxRetries:     5,
		promHistograms: false,
		extract:        false,
		secret:         "",
		security:       false,
		AuthFunc:       nil,
		schemas:        []string{"basic", "bearer"},
	}

	var err error
	var pool *x509.CertPool
	var cert tlspkg.Certificate
	var opts []grpc.ServerOption

	addr := ip
	if len(ip) == 0 {
		ip = lib.ResolveHostIp()
	} else if ip == "0" {
		ip = lib.ResolveHostIp()
		addr = ""
	}
	if port == 0 {
		port, err = lib.GetPort(ip)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	server.ip = ip
	server.port = port
	server.Addr = fmt.Sprintf("%s:%d", addr, port)
	server.ServiceID = fmt.Sprintf("%s:%s:%d", name, ip, port)

	for _, option := range options {
		option(server)
	}

	log.Debugln("debug:", server.debug, "verbose:", server.verbose, "ServiceID:", server.ServiceID, "port:", server.port)
	// Server options

	if server.tls {
		// certFile = flag.String("cert", "", "Certificate file")
		// keyFile  = flag.String("key", "", "Key file")
		cert, err = tlspkg.LoadX509KeyPair(server.certFile, server.keyFile)
		if err != nil {
			log.Fatalln("Cannot load certificate", err.Error())
		}
		// Create pool to trust
		caCert, err := ioutil.ReadFile(server.certFile)
		if err != nil {
			log.Fatalln("Cannot load certificate", err.Error())
		}
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
		// We don't need the instruct gRPC to do TLS because we are using cmux to proxy TLS
		creds := credentials.NewServerTLSFromCert(&cert)
		opts = append(opts, grpc.Creds(creds))
	}

	// RateLimiter - Выполняется при приходе запроса на соединение
	// Это логика до того как мы начнем чтото парсить
	// qps      = flag.Float64("qps", 5, "Queries per second in rate limiter")
	// burst    = flag.Int("burst", 1, "Burst in rate limiter")
	// tap := grpcserver.NewTapHandler(
	//   grpcserver.NewMetrics(),
	//   // rate.Limit(*qps),
	//   rate.Limit(5),
	//   // *burst,
	//   1,
	// )

	// Common options
	// opts = append(opts, grpc.MaxRecvMsgSize(1<<20)) // 1MB
	// opts = append(opts, grpc.InTapHandle(tap.Handle))

	var optStreams []grpc.StreamServerInterceptor
	var optUnarys []grpc.UnaryServerInterceptor

	if server.isLogrus {
		if server.verbose > 1 {
			log.Println("Logrus logger enabled...")
		}
		optLogger := []grpclogrus.Option{
			grpclogrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
				return "grpc.time_ns", duration.Nanoseconds()
			}),
		}

		// Make sure that log statements internal to gRPC library are logged using the logrus Logger as well.
		grpclogrus.ReplaceGrpcLogger(server.LogrusEntry)

		if server.extract {
			if server.verbose > 1 {
				log.Println("CodeGenRequestFieldExtractor enabled...")
			}
			optStreams = append(optStreams, ctxtags.StreamServerInterceptor(ctxtags.WithFieldExtractor(ctxtags.CodeGenRequestFieldExtractor)))
			optUnarys = append(optUnarys, ctxtags.UnaryServerInterceptor(ctxtags.WithFieldExtractor(ctxtags.CodeGenRequestFieldExtractor)))
		} else {
			optStreams = append(optStreams, ctxtags.StreamServerInterceptor())
			optUnarys = append(optUnarys, ctxtags.UnaryServerInterceptor())
		}

		optStreams = append(optStreams, grpclogrus.StreamServerInterceptor(server.LogrusEntry, optLogger...))
		optUnarys = append(optUnarys, grpclogrus.UnaryServerInterceptor(server.LogrusEntry, optLogger...))
	}

	if server.prom {
		if server.verbose > 1 {
			log.Println("Prometheus enabled...")
		}
		if server.promHistograms {
			grpcprom.EnableHandlingTimeHistogram()
		}
		optStreams = append(optStreams, grpcprom.StreamServerInterceptor)
		optUnarys = append(optUnarys, grpcprom.UnaryServerInterceptor)
	}
	// if server.openTracing {
	// if server.verbose {
	//  log.Println("Open Tracing enabled...")
	// }
	// 	optStreams = append(optStreams, grpcopentracing.StreamServerInterceptor())
	// 	optUnarys = append(optUnarys, grpcopentracing.UnaryServerInterceptor())
	// }

	if server.validation {
		if server.verbose > 1 {
			log.Println("Validation enabled...")
		}
		optStreams = append(optStreams, grpcvalidtor.StreamServerInterceptor())
		optUnarys = append(optUnarys, grpcvalidtor.UnaryServerInterceptor())
	}

	if server.debug {
		log.Println("Debug enabled...")
		optUnarys = append(optUnarys, TimingServerInterceptor)
	}

	if server.isAuth {
		if server.verbose > 1 {
			log.Println("Authorization enabled...")
		}
		if server.AuthFunc == nil {
			optStreams = append(optStreams, grpcauth.StreamServerInterceptor(server.ServiceAuthFunc))
			optUnarys = append(optUnarys, grpcauth.UnaryServerInterceptor(server.ServiceAuthFunc))
		} else {
			log.Warningln("TODO Authorization custom function")
			optStreams = append(optStreams, grpcauth.StreamServerInterceptor(Authenticate))
			optUnarys = append(optUnarys, grpcauth.UnaryServerInterceptor(Authenticate))
		}
	}

	if len(server.middlewareUnarys) > 0 {
		if server.verbose > 1 {
			log.Println("custom middleware")
		}
		for _, m := range server.middlewareUnarys {
			optUnarys = append(optUnarys, m)
		}
	}
	if len(server.middlewareStreams) > 0 {
		if server.verbose > 1 {
			log.Println("custom middleware stream")
		}
		for _, m := range server.middlewareStreams {
			optStreams = append(optStreams, m)
		}
	}

	if server.recovery {
		if server.verbose > 1 {
			log.Println("Recovery enabled...")
		}
		optStreams = append(optStreams, grpcrecovery.StreamServerInterceptor())
		optUnarys = append(optUnarys, grpcrecovery.UnaryServerInterceptor())
	}

	// gRPC middleware
	if len(optStreams) > 0 {
		opts = append(opts, grpc.StreamInterceptor(grpcmw.ChainStreamServer(optStreams...)))
	}
	if len(optUnarys) > 0 {
		opts = append(opts, grpc.UnaryInterceptor(grpcmw.ChainUnaryServer(optUnarys...)))
	}

	server.Conn = grpc.NewServer(opts...)

	if server.prom {
		if server.verbose > 1 {
			log.Println("Register Prometheus")
		}
		grpcprom.Register(server.Conn)
	}

	return server
}

// Serve accepts incoming connections
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s *Server) Serve(l net.Listener) error {
	s.Register(false)
	return s.Conn.Serve(l)
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// If srv.Addr is blank, ":http" is used.
// ListenAndServe always returns a non-nil error.
func (s *Server) ListenAndServe() error {
	list, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		if s.verbose > 2 {
			log.Printf("start server: %v\n", s.Addr)
		}
		s.Register(true)
	}
	s.DeferShutdown()
	return s.Conn.Serve(list)
}

func (s *Server) DeferShutdown() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		x := <-ch
		if s.verbose > 2 {
			log.Println("[server] receive signal: ", x)
		}
		s.GracefulStop()
		if s.verbose > 0 {
			log.Println("STOP gracefull gRPC server:", s.ServiceID)
		}
	}()
}

// Stop stops the gRPC server. It immediately closes all open
// connections and listeners.
// It cancels all active RPCs on the server side and the corresponding
// pending RPCs on the client side will get notified by connection
// errors.
func (s *Server) Stop() {
	s.Conn.Stop()
}

// GracefulStop stops the gRPC server gracefully. It stops the server from
// accepting new connections and RPCs and blocks until all the pending RPCs are
// finished.
func (s *Server) GracefulStop() {
	if s.verbose > 2 {
		log.Println("----- GracefulStop gRPC ----")
	}
	s.Conn.GracefulStop()
}

// GetPort порт gRPC сервера
func (s *Server) GetPort() int {
	return s.port
}

// GetIP - внешний IP gRPC сервера
func (s *Server) GetIP() string {
	return s.ip
}

// Register is the helper function to self-register service into Consul server
// target - consul dial address, for example: "127.0.0.1:8500"
// interval - interval of self-register to etcd
// ttl - ttl of the register information
func (s *Server) Register(unregister bool) {
	// s.ConsulAddress = target
	if s.isConsul {
		if err := cr.Register(s.Name, s.Namespace, s.ip, s.port, s.ConsulAddress, s.consulInterval, s.consulTtl, unregister, s.verbose); err != nil {
			log.Println(err.Error())
		}
	}
}

// UnRegister is the helper function to unregister service into Consul server
func (s *Server) UnRegister() {
	if s.isConsul {
		cr.UnRegister(s.Name, s.ip, s.port, s.ConsulAddress, s.verbose)
	}
}

func SetDomain(domain string) ServerOption {
	return func(s *Server) {
		s.Domain = domain
	}
}

func SetDescription(description string) ServerOption {
	return func(s *Server) {
		s.Description = description
	}
}

// SetSetting разные настройки
func SetSetting(verbose int, debug bool, prom bool) ServerOption {
	return func(s *Server) {
		s.debug = debug
		s.verbose = verbose
		// s.openTracing = openTracing
		s.prom = prom
	}
}

func SetConsul(addr string, interval time.Duration, ttl int) ServerOption {
	return func(s *Server) {
		if addr != "none" {
			s.isConsul = true
			s.consulInterval = interval
			s.consulTtl = ttl
			s.ConsulAddress = addr
		} else {
		  s.isConsul = false
    }
	}
}

// Use - custom Unary middleware
func Use(item grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.middlewareUnarys = append(s.middlewareUnarys, item)
	}
}

// UseStream - custom stream middleware
func UseStream(item grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.middlewareStreams = append(s.middlewareStreams, item)
	}
}

// DisbaleRecovery отключаем ывостановление
func DisbaleRecovery() ServerOption {
	return func(s *Server) {
		s.recovery = false
	}
}

// DisbaleValidation отключаем проверки
func DisbaleValidation() ServerOption {
	return func(s *Server) {
		s.validation = false
	}
}

// SetExtract - CodeGenRequestFieldExtractor is a function that relies on code-generated functions
// that export log fields from requests.
// These are usually coming from a protoc-plugin that generates additional information based on custom field options.
func SetExtract() ServerOption {
	return func(s *Server) {
		s.extract = true
	}
}

// SetLogger устанавливаем значения для Logrus
func SetLogger(logger *log.Logger) ServerOption {
	return func(s *Server) {
		s.isLogrus = true
		s.LogrusEntry = log.NewEntry(logger)
		// s.LogLevel = level
		// s.Logger = logger
	}
}

// SetAuth устанавливаем значения для авторизации
func SetAuth(service *authorize.Authorize) ServerOption {
	return func(s *Server) {
		s.isAuth = true
		s.secret = service.Secret
		s.security = service.Security
		s.Authorize = service
	}
}

// SetCritical критичекое количество сервисов
//  0 - некритично
//  число работающих екземпляров сервиса  например 1 если сервисов будет меньше одного то произойдет событие
func SetCritical(val uint) ServerOption {
	return func(s *Server) {
		s.Critical = val
	}
}

// ServiceAuthFunc authenticated function
func (s *Server) ServiceAuthFunc(ctx context.Context) (context.Context, error) {
	cnt := len(s.schemas)
	var token string
	var err error
	for _, schema := range s.schemas {
		token, err = grpcauth.AuthFromMD(ctx, schema)
		if err != nil {
			if cnt == 1 {
				log.Errorln(err.Error())
				return nil, err
			}
		} else {
			if schema == "bearer" {
				tokenInfo, err := s.Authorize.ParseToken(token)
				if err != nil {
					log.Errorf("invalid auth token: %v", err.Error())
					return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err.Error())
				}
				user := "user"
				if len(tokenInfo.Prefix) > 0 {
					user = fmt.Sprintf("%s.%s", tokenInfo.Prefix, user)
				}
				ctxtags.Extract(ctx).Set(user, tokenInfo.User)

				userId := "user_id"
				if len(tokenInfo.Prefix) > 0 {
					userId = fmt.Sprintf("%s.%s", tokenInfo.Prefix, userId)
				}
				ctxtags.Extract(ctx).Set(userId, tokenInfo.UserId)

				for _, param := range tokenInfo.Params {
					name := param.Name
					if len(tokenInfo.Prefix) > 0 {
						ctxtags.Extract(ctx).Set(name, param.Value)
					}
				}

				newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
				return newCtx, nil
			} else {
				tokenInfo, err := s.Authorize.ParseBasic(token)
				if err != nil {
					log.Errorf("[Basic Auth] %v", err.Error())
					return nil, status.Errorf(codes.Unauthenticated, "[Basic Auth] %v", err.Error())
				}
				user := "user"
				if len(tokenInfo.Prefix) > 0 {
					user = fmt.Sprintf("%s.%s", tokenInfo.Prefix, user)
				}
				ctxtags.Extract(ctx).Set(user, tokenInfo.User)

				userId := "user_id"
				if len(tokenInfo.Prefix) > 0 {
					userId = fmt.Sprintf("%s.%s", tokenInfo.Prefix, userId)
				}
				ctxtags.Extract(ctx).Set(userId, tokenInfo.UserId)

				newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
				return newCtx, nil
			}
		}
	}
	log.Errorln(err.Error())
	return nil, status.Errorf(codes.Unauthenticated, "is not support auth schema")
}
