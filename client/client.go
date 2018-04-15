package client

import (
	"crypto/x509"
	"fmt"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/micro-grpc/mbox/lb"
	"github.com/micro-grpc/mbox/lib"
	"github.com/olivere/grpc/lb/healthz"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/naming"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// "google.golang.org/grpc/balancer/roundrobin"
// grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
// grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"

// Client
type Client struct {
	Name           string // имя сервиса для Consul
	Namespace      string
	Conn           *grpc.ClientConn
	tokenAuth      *TokenAuth
	Addr           string
	HealthChecks   []string
	isConsul       bool
	isAuth         bool
	isLogrus       bool
	tls            bool
	ServerName     string // "" если не использовать TLS
	ServiceName    string // имя сервиса для Consul
	ConsulAddress  string
	caFile         string
	limiter        *rate.Limiter
	maxRetries     uint
	Balancer       grpc.Balancer
	ConsulResolver *lb.ConsulResolver
	verbose        int
	debug          bool
	prom           bool
	promHistograms bool
	LogrusEntry    *log.Entry
	schema         string
	token          string
}

type ClientOption func(*Client)

func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		Namespace:      "",
		Addr:           "localhost:9000",
		HealthChecks:   nil,
		isConsul:       false,
		isAuth:         false,
		isLogrus:       false,
		tls:            false,
		ServerName:     "",
		ServiceName:    "",
		ConsulAddress:  "",
		caFile:         "",
		limiter:        rate.NewLimiter(rate.Limit(1000), 10),
		maxRetries:     5,
		verbose:        0,
		debug:          false,
		prom:           false,
		promHistograms: false,
		schema:         "none",
		token:          "",
	}

	for _, option := range options {
		option(client)
	}

	var opts []grpc.DialOption

	// Configure TLS
	if client.tls {
		cert, err := ioutil.ReadFile(client.caFile)
		if err != nil {
			return nil, errors.Wrap(err, "cannot read caFile")
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(cert) {
			return nil, errors.New("failed to append certificate to pool")
		}
		var sn string
		if client.ServerName != "" {
			sn = client.ServerName
		} else {
			sn, _, err = net.SplitHostPort(client.Addr)
			if err != nil {
				return nil, errors.Wrap(err, "cannot split address into host and port")
			}
		}
		creds := credentials.NewClientTLSFromCert(pool, sn)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if client.isAuth {
		opts = append(opts, grpc.WithPerRPCCredentials(client.tokenAuth))
	}

	var optStreams []grpc.StreamClientInterceptor
	var optUnarys []grpc.UnaryClientInterceptor

	if client.isLogrus {
		if client.verbose > 0 {
			log.Println("Logrus logger enabled...")
		}
		optLogger := []grpclogrus.Option{
			grpclogrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
				return "grpc.time_ns", duration.Nanoseconds()
			}),
		}

		optStreams = append(optStreams, grpclogrus.StreamClientInterceptor(client.LogrusEntry, optLogger...))
		optUnarys = append(optUnarys, grpclogrus.UnaryClientInterceptor(client.LogrusEntry, optLogger...))
	}

	// Monitoring via Prometheus
	if client.prom {
		if client.verbose > 0 {
			log.Println("Prometheus enabled...")

		}
		optStreams = append(optStreams, grpcprom.StreamClientInterceptor)
		optUnarys = append(optUnarys, grpcprom.UnaryClientInterceptor)
	}

	// Retries
	if client.verbose > 0 {
		log.Println("Retries enabled...")
	}
	retrycallopts := []grpcretry.CallOption{
		grpcretry.WithMax(client.maxRetries),
		grpcretry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
	}
	optStreams = append(optStreams, grpcretry.StreamClientInterceptor(retrycallopts...))
	optUnarys = append(optUnarys, grpcretry.UnaryClientInterceptor(retrycallopts...))

	opts = append(opts, grpc.WithStreamInterceptor(grpcmw.ChainStreamClient(optStreams...)))
	opts = append(opts, grpc.WithUnaryInterceptor(grpcmw.ChainUnaryClient(optUnarys...)))

	opts = append(opts, grpc.WithBlock())

	if client.isConsul {
		if client.verbose > 0 {
			log.Println("consul enabled...")
		}
		client.ConsulResolver = lb.NewResolver(client.ServiceName, client.Namespace, client.ConsulAddress, client.Addr)
		client.Balancer = grpc.RoundRobin(client.ConsulResolver)
		client.Balancer.Start(client.Addr, grpc.BalancerConfig{})
		//
		opts = append(opts, grpc.WithBalancerName(roundrobin.Name))
		opts = append(opts, grpc.WithBalancer(client.Balancer))
		// client.Addr = client.ConsulAddress
	}

	var err error

	if client.verbose > 0 {
		log.Printf("Dial: %s ServiceName: %s namespace: %s\n", client.Addr, client.ServiceName, client.Namespace)
	}
	client.Conn, err = grpc.Dial(client.Addr, opts...)

	return client, err
}

func (c *Client) Close() error {
	return c.Conn.Close()
}

// HealthzResolver - пока непонятно
func (c *Client) HealthzResolver() (naming.Resolver, error) {
	if len(c.HealthChecks) == 0 {
		return nil, errors.New("no healthcheck URLs specified")
	}
	addrs := strings.Split(c.Addr, ",")
	if want, have := len(addrs), len(c.HealthChecks); want != have {
		return nil, errors.Errorf("there must be a healthcheck URL for every gRPC endpoint; "+
			"you passed %d gRPC endpoints but have %d healthcheck URLs", want, have)
	}

	var endpoints []healthz.Endpoint
	for i, addr := range addrs {
		healthcheckURL := c.HealthChecks[i]
		if _, err := url.Parse(healthcheckURL); err != nil {
			return nil, errors.Wrapf(err, "invalid URL: %s", healthcheckURL)
		}
		endpoints = append(endpoints, healthz.Endpoint{
			Addr:     addr,
			CheckURL: healthcheckURL,
		})
	}

	r, err := healthz.NewResolver(healthz.SetEndpoints(endpoints...))
	if err != nil {
		return nil, errors.Wrap(err, "error creating healthz resolver")
	}
	return r, nil
}

func SetAddr(addr string) ClientOption {
	return func(client *Client) {
		client.Addr = addr
	}
}

func SetHealthcheckURL(urls ...string) ClientOption {
	return func(client *Client) {
		client.HealthChecks = urls
	}
}

func SetTLS(tls bool) ClientOption {
	return func(client *Client) {
		client.tls = tls
	}
}

func SetServerName(serverName string) ClientOption {
	return func(client *Client) {
		client.ServerName = serverName
	}
}

func SetCAFile(caFile string) ClientOption {
	return func(client *Client) {
		client.caFile = caFile
	}
}

func SetRateLimiter(limiter *rate.Limiter) ClientOption {
	return func(client *Client) {
		client.limiter = limiter
	}
}

func SetMaxRetries(maxRetries uint) ClientOption {
	return func(client *Client) {
		client.maxRetries = maxRetries
	}
}

func SetConsul(isConsul bool, consulAddr string) ClientOption {
	return func(client *Client) {
		if isConsul {
			client.ConsulAddress = consulAddr
			client.isConsul = true
		}
	}
}

// SetSetting разные настройки
func SetSetting(verbose int, debug bool, prom bool) ClientOption {
	return func(client *Client) {
		client.debug = debug
		client.verbose = verbose
		// s.openTracing = openTracing
		client.prom = prom
	}
}

// SetServiceName созаем имя сервиса и его адресс
// name - имя сервиса разделитель - для того чтобы можно было к нему обращатся через DNS
// addr адрес сервера к которому будет конектится клиент
//      если "" то будет автоматически братся IP внешнего интерфейса
//      если local то будет коннектится к localhost:9000
//      port - порт на котором прослушивает сервер
func SetServiceName(name string, namespace string, addr string, port int) ClientOption {
	return func(client *Client) {
		client.ServiceName = name
		client.Namespace = namespace
		if len(addr) == 0 {
			ip := lib.ResolveHostIp()
			client.Addr = fmt.Sprintf("%s:%d", ip, port)
		} else if addr != "local" {
			client.Addr = addr
		}
	}
}

// SetAuthJWT - setting
func SetAuthJWT(token string) ClientOption {
	return func(client *Client) {
		client.isAuth = true
		client.schema = "bearer"
		client.token = token
		client.tokenAuth = &TokenAuth{Token: token, Security: false, Schema: "Bearer"}
	}
}

// SetAuthBasic - setting Basic Authentication
// provided login and passwd
func SetAuthBasic(login string, passwd string) ClientOption {
	return func(client *Client) {
		client.isAuth = true
		client.schema = "basic"
		client.token = lib.BasicAuth(login, passwd)
		client.tokenAuth = &TokenAuth{Token: client.token, Security: false, Schema: "Basic"}
		// client.tokenAuth = &TokenAuth{Token: client.token, Security: false}
		// if isAuth {
		//   client.tokenAuth = &TokenAuth{Token: token, Security: security}
		// }
	}
}

// SetLogger устанавливаем значения для Logrus
func SetLogger(logger *log.Logger) ClientOption {
	return func(client *Client) {
		client.isLogrus = true
		client.LogrusEntry = log.NewEntry(logger)
		// s.LogLevel = level
		// s.Logger = logger
	}
}

func SetBalancer(balancer grpc.Balancer) ClientOption {
	return func(client *Client) {
		client.Balancer = balancer
	}
}

func (c *Client) Shutdown() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		x := <-ch
		if c.verbose > 0 {
			log.Println("[client] micro-grpc: receive signal: ", x)
		}
	}()
}

func (c *Client) GetServerAddress() (addr string) {
	if c.isConsul {
		return c.ConsulResolver.GetFirst()
	}
	return c.Addr
}
