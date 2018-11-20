package consul

import (
  "fmt"
  consulapi "github.com/hashicorp/consul/api"
  log "github.com/sirupsen/logrus"
)

// GetKV - получить данные для клюса
func GetKV(prefix string, user string, key string, defaultVal string, conn *consulapi.Client) (val string) {
  k := fmt.Sprintf("%s/%s/%s", prefix, user, key)
  if len(key) == 0 {
    k = fmt.Sprintf("%s/%s", prefix, user)
  }
  kvp, _, err := conn.KV().Get(k, nil)

  if err != nil {
    log.Errorln(err.Error())
    return defaultVal
  }
  if  kvp == nil {
    return defaultVal
  }
  return string(kvp.Value)
}

// AddKV - добавляем ключ в Consul KV
func AddKV(prefix string, user string, key string, val string, conn *consulapi.Client) (err error, mess string) {
  k := fmt.Sprintf("%s/%s/%s", prefix, user, key)
  if len(key) == 0 {
    k = fmt.Sprintf("%s/%s", prefix, user)
  }
  item := &consulapi.KVPair{Key: k, Value: []byte(val)}
  // Consul.KV().Acquire(item, nil)
  res, err := conn.KV().Put(item, nil)
  if err != nil {
    //panic(err)
    return err, ""
  }
  return nil, fmt.Sprintf("затрачено: %v", res.RequestTime)
}

type ClientConsul struct {
  Client  *consulapi.Client
  Address string
  Config  *consulapi.Config
  Verbose int
  Err     error
}

func NewClientConsul(addr string, verbose int) *ClientConsul {
  configConsul := consulapi.DefaultConfig()
  configConsul.Address = addr

  c := ClientConsul{
    Address: addr,
    Config:  configConsul,
    Verbose: verbose,
  }

  c.Client, c.Err = consulapi.NewClient(configConsul)
  return &c
}

func (c *ClientConsul) GetKVstring(path string, defaultVal string) (val string, err error) {
  kvPair, meta, err := c.Client.KV().Get(path, nil)
  if err != nil {
    return defaultVal, err
  }
  if c.Verbose > 0 {
    log.Println("[CONSUL:GET] response time:", meta.RequestTime)
  }
  if kvPair == nil {
    return defaultVal, nil
  }
  return string(kvPair.Value), nil
}

func (c *ClientConsul) GetKV(path string) (val []byte, err error) {
  kvPair, meta, err := c.Client.KV().Get(path, nil)
  if err != nil {
    return nil, err
  }
  if c.Verbose > 0 {
    log.Println("[CONSUL:GET] response time:", meta.RequestTime)
  }
  if kvPair == nil {
    return nil, nil
  }
  return kvPair.Value, nil
}

func (c *ClientConsul) AddKVstring(key string, val string) (err error) {
  item := &consulapi.KVPair{Key: key, Value: []byte(val)}
  res, err := c.Client.KV().Put(item, nil)

  if err != nil {
    //panic(err)
    return err
  }

  if c.Verbose > 0 {
    log.Println("[CONSUL:PUT] response time:", res.RequestTime)
  }
  return nil
}

func (c *ClientConsul) AddKV(key string, val []byte) (err error) {
  item := &consulapi.KVPair{Key: key, Value: val}
  res, err := c.Client.KV().Put(item, nil)

  if err != nil {
    //panic(err)
    return err
  }

  if c.Verbose > 0 {
    log.Println("[CONSUL:PUT] response time:", res.RequestTime)
  }
  return nil
}
