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
	sudoUsageMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "infra_sudo_usage_total",
			Help: "총 sudo 명령 사용 횟수",
		},
	)

	sudoLogPath     = "/var/log/auth.log" // Ubuntu 기준
	sudoLastOffset  int64
	sudoUsageMutex  sync.Mutex
)

func RegisterSudoUsageCollector() {
	prometheus.MustRegister(sudoUsageMetric)

	go func() {
		for {
			parseSudoUsage()
			time.Sleep(10 * time.Second)
		}
	}()
}

func parseSudoUsage() {
	sudoUsageMutex.Lock()
	defer sudoUsageMutex.Unlock()

	file, err := os.Open(sudoLogPath)
	if err != nil {
		log.Printf("Sudo 로그 파일 열기 실패: %v", err)
		return
	}
	defer file.Close()

	// 이전 읽기 위치로 이동
	if sudoLastOffset > 0 {
		_, err := file.Seek(sudoLastOffset, 0)
		if err != nil {
			log.Printf("파일 seek 실패: %v", err)
			return
		}
	}

	scanner := bufio.NewScanner(file)
	sudoRegex := regexp.MustCompile(`sudo: .*: TTY=`)
	var count int

	for scanner.Scan() {
		line := scanner.Text()
		if sudoRegex.MatchString(line) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Sudo 로그 스캔 실패: %v", err)
		return
	}

	if count > 0 {
		sudoUsageMetric.Add(float64(count))
	}

	// 현재 읽은 위치 저장
	offset, err := file.Seek(0, os.SEEK_CUR)
	if err == nil {
		sudoLastOffset = offset
	}
}

