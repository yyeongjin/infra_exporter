package collector

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	listeningPorts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "infra_listening_ports",
			Help: "리스닝 중인 포트의 개수",
		},
		[]string{"port", "protocol"},
	)

	portMutex sync.Mutex
)

func RegisterPortCollector() {
	prometheus.MustRegister(listeningPorts)

	go func() {
		for {
			collectListeningPorts()
			time.Sleep(1 * time.Minute)
		}
	}()
}

func collectListeningPorts() {
	portMutex.Lock()
	defer portMutex.Unlock()

	listeningPorts.Reset()

	cmd := exec.Command("netstat", "-tuln")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("netstat 명령 실행 실패: %v", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "tcp") || strings.HasPrefix(line, "udp") {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			protocol := fields[0]
			address := fields[3]
			port := parsePort(address)
			if port != "" {
				listeningPorts.WithLabelValues(port, protocol).Set(1)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("포트 정보 스캔 실패: %v", err)
	}
}

func parsePort(address string) string {
	// 예시: 0.0.0.0:22 또는 [::]:443
	if strings.HasPrefix(address, "[") {
		// IPv6
		parts := strings.Split(address, "]:")
		if len(parts) == 2 {
			return parts[1]
		}
	} else {
		// IPv4
		parts := strings.Split(address, ":")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}
	return ""
}

