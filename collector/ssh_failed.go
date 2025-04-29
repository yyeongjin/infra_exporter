package collector

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	sshFailedMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "infra_ssh_failed_logins_total",
			Help: "총 SSH 로그인 실패 횟수",
		},
	)

	sshFailedMutex sync.Mutex
	sshLogPath     = "/var/log/auth.log"
	lastOffset     int64
)

func RegisterSSHFailedCollector() {
	prometheus.MustRegister(sshFailedMetric)

	go func() {
		for {
			parseSSHFailed()
			time.Sleep(10 * time.Second)
		}
	}()
}

func parseSSHFailed() {
	sshFailedMutex.Lock()
	defer sshFailedMutex.Unlock()

	file, err := os.Open(sshLogPath)
	if err != nil {
		log.Printf("SSH 로그 파일 열기 실패: %v", err)
		return
	}
	defer file.Close()

	// 이전 읽기 위치로 이동
	if lastOffset > 0 {
		_, err := file.Seek(lastOffset, 0)
		if err != nil {
			log.Printf("파일 seek 실패: %v", err)
			return
		}
	}

	scanner := bufio.NewScanner(file)
	failRegex := regexp.MustCompile(`Failed password for`)
	var count int

	for scanner.Scan() {
		line := scanner.Text()
		if failRegex.MatchString(line) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("SSH 로그 스캔 실패: %v", err)
		return
	}

	if count > 0 {
		sshFailedMetric.Add(float64(count))
	}

	// 현재 위치 저장
	offset, err := file.Seek(0, os.SEEK_CUR)
	if err == nil {
		lastOffset = offset
	}
}

