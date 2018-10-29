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
