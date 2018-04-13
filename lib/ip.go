package lib

import (
	"fmt"
	"net"
)

// ResolveHostIp - IP Адресс внешнего интерфейса
func ResolveHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ip := networkIp.IP.String()
			// fmt.Println("Resolved Host IP: " + ip)
			return ip
		}
	}
	return ""
}

func CheckFreePort(ip string, port int) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetFreePort - поиск свободного порта из диапазона
//  port, err := core.GetFreePort("localhost", 9000, 9999)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println("free port: ", port)
func GetFreePort(ip string, start int, end int) (port int, err error) {
	for start < end {
		port, err := CheckFreePort(ip, start)
		if err == nil {
			return port, nil
		}
		start++
	}
	return 0, fmt.Errorf("Is NOT free port start: %d end: %d\n", start, end)
}

// GetPort свободный порт в диапазоне от 9000 до 50000
func GetPort(ip string) (port int, err error) {
	startPort := 9000
	endPort := startPort + 50000
	port, err = GetFreePort(ip, startPort, endPort)
	return port, err
}
