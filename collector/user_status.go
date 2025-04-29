package collector

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	userStatusMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "infra_user_status",
			Help: "사용자 계정 상태 (1: 활성, 0: 비활성)",
		},
		[]string{"user"},
	)

	userStatusFilePath = "/etc/passwd"
	userStatusMutex    sync.Mutex
)

func RegisterUserStatusCollector() {
	prometheus.MustRegister(userStatusMetric)

	go func() {
		for {
			checkUserStatus()
			time.Sleep(60 * time.Second)
		}
	}()
}

func checkUserStatus() {
	userStatusMutex.Lock()
	defer userStatusMutex.Unlock()

	file, err := os.Open(userStatusFilePath)
	if err != nil {
		log.Printf("사용자 상태 파일 열기 실패: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue // 필드가 부족한 경우 skip
		}

		user := fields[0]
		shell := fields[6]

		// 쉘이 /sbin/nologin 또는 /bin/false인 경우 비활성
		if shell == "/sbin/nologin" || shell == "/bin/false" {
			userStatusMetric.WithLabelValues(user).Set(0)
		} else {
			userStatusMetric.WithLabelValues(user).Set(1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("사용자 상태 파일 스캔 실패: %v", err)
	}
}

