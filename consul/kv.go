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
