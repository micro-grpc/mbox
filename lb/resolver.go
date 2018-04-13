package lb

import (
  "errors"
  "fmt"

  consulapi "github.com/hashicorp/consul/api"
  "google.golang.org/grpc/naming"
)

// ConsulResolver is the implementaion of grpc.naming.Resolver
type ConsulResolver struct {
  ServiceName string // service name
  Target string
  Addr string
  Addresses []string
  Scheme string
  Verbose int
  Namespace string
}

// GetFirst - возвращает первый адрес из списка
func (cr *ConsulResolver) GetFirst() (addr string) {
  if len(cr.Addresses) > 0 {
    return cr.Addresses[0]
  }
  return cr.Addr
}

// NewResolver return ConsulResolver with service name
func NewResolver(serviceName string, namespace string, target string, addr string) *ConsulResolver {
  return &ConsulResolver{ServiceName: serviceName, Namespace: namespace, Target: target, Addr: addr, Verbose: 0}
}

// Resolve to resolve the service from consul, target is the dial address of consul
func (cr *ConsulResolver) Resolve(target string) (naming.Watcher, error) {
  if cr.ServiceName == "" {
    return nil, errors.New("no service name provided")
  }

  if cr.Verbose > 1 {
    fmt.Println("ServiceName:", cr.ServiceName, "target:", target, "consul:", cr.Target)
  }
  conf := consulapi.DefaultConfig()
  conf.Address = cr.Target
  // generate consul client, return if error
  // conf := &consul.Config{
  //   Scheme:  "http",
  //   Address: target,
  // }
  client, err := consulapi.NewClient(conf)
  if err != nil {
    return nil, fmt.Errorf("consul conect error: %v", err)
  }

  // return ConsulWatcher
  watcher := &ConsulWatcher{
    cr: cr,
    cc: client,
  }
  return watcher, nil
}
