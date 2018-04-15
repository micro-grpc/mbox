package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func UnRegister(name string, host string, port int, target string, verbose int) {
	serviceID := fmt.Sprintf("%s:%s:%d", name, host, port)

	if verbose > 2 {
		log.Println("[START] UnRegister", serviceID)
	}

	conf := consulapi.DefaultConfig()
	conf.Address = target

	client, err := consulapi.NewClient(conf)
	if err != nil {
		log.Errorf("[UnRegister] create consul client error: %v", err.Error())
	}

	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Errorf("[UnRegister] deregister service '%s' error: %v\n", serviceID, err.Error())
	} else if verbose > 0 {
		log.Printf("[UnRegister] deregistered service '%s' from consul server.\n", serviceID)
	}

	if err := client.Agent().CheckDeregister(serviceID); err != nil {
		log.Errorf("[UnRegister] deregister check '%s' error: %v\n", serviceID, err.Error())
	}
}

// Register is the helper function to self-register service into Etcd/Consul server
// name - service name
// host - service host
// port - service port
// target - consul dial address, for example: "127.0.0.1:8500"
// interval - interval of self-register to etcd
// ttl - ttl of the register information
// debug - show register message
func Register(name string, namespace string, host string, port int, target string, interval time.Duration, ttl int, unregister bool, verbose int) error {

	// conf := &consulapi.Config{Scheme: "http", Address: target}
	conf := consulapi.DefaultConfig()
	conf.Address = target

	client, err := consulapi.NewClient(conf)
	if err != nil {
		return fmt.Errorf("[Register] create consul client error: %v", err)
	}

	serviceID := fmt.Sprintf("%s:%s:%d", name, host, port)

	// de-register if meet signhup
	if unregister {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
			x := <-ch
			if verbose > 2 {
				log.Println("[consul] micro-grpc: receive signal: ", x)
			}

			err := client.Agent().ServiceDeregister(serviceID)
			if err != nil {
				log.Errorf("[Register] deregister service '%s' error: %v\n", serviceID, err.Error())
			} else if verbose > 0 {
				log.Printf("[Register] deregistered service '%s' from consul server.\n", serviceID)
			}

			if err := client.Agent().CheckDeregister(serviceID); err != nil {
				log.Errorf("[Register] deregister check '%s' error: %v\n", serviceID, err.Error())
			}
			// s, _ := strconv.Atoi(fmt.Sprintf("%d", x))
			// os.Exit(s)
		}()
	}

	// routine to update ttl
	go func() {
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			err := client.Agent().UpdateTTL(serviceID, "", "passing")
			if err != nil {
				log.Errorln("[Register] update ttl of service error: ", err.Error())
			}
		}
	}()

	// initial register service
	regis := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Address: host,
		Port:    port,
	}
	if len(namespace) > 0 {
		regis.Tags = []string{namespace}
	}
	err = client.Agent().ServiceRegister(regis)
	if err != nil {
		return fmt.Errorf("[Register] initial register service '%s' host to consul error: %s", name, err.Error())
	}

	// initial register service check
	check := consulapi.AgentServiceCheck{TTL: fmt.Sprintf("%ds", ttl), Status: "passing"}
	err = client.Agent().CheckRegister(&consulapi.AgentCheckRegistration{
		ID:                serviceID,
		Name:              name,
		ServiceID:         serviceID,
		AgentServiceCheck: check,
	})
	if err != nil {
		return fmt.Errorf("[Register] initial register service check to consul error: %s", err.Error())
	}

	if verbose > 1 {
		log.Printf("Register service '%s' OK!\n", serviceID)
		if len(namespace) > 0 {
			log.Printf("Register service '%s' tag: %s OK!\n", serviceID, namespace)
		}
	}
	return nil
}
