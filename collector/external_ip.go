package collector

import (
	"log"
	"net"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	externalIPMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "infra_external_ip_connection_total",
			Help: "외부 IP 연결 빈도 (상위 N개)",
		},
		[]string{"ip"},
	)

	externalIPMutex sync.Mutex
)

func RegisterExternalIPCollector(topN int) {
	prometheus.MustRegister(externalIPMetric)

	go func() {
		for {
			collectExternalIP(topN)
			time.Sleep(1 * time.Minute)
		}
	}()
}

func collectExternalIP(topN int) {
	externalIPMutex.Lock()
	defer externalIPMutex.Unlock()

	// netstat 명령어 실행
	out, err := exec.Command("netstat", "-ntu").Output()
	if err != nil {
		log.Printf("외부 IP 수집 실패: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	ipCount := make(map[string]int)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		if fields[0] != "tcp" && fields[0] != "udp" {
			continue
		}
		remoteAddr := fields[4]
		host, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			continue
		}
		if isPrivateIP(host) {
			continue
		}
		ipCount[host]++
	}

	// 기존 메트릭 초기화
	externalIPMetric.Reset()

	// IP 빈도수 정렬
	type ipFreq struct {
		ip    string
		count int
	}
	var freqList []ipFreq
	for ip, count := range ipCount {
		freqList = append(freqList, ipFreq{ip, count})
	}
	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].count > freqList[j].count
	})

	// 상위 N개 메트릭 등록
	for i := 0; i < len(freqList) && i < topN; i++ {
		externalIPMetric.WithLabelValues(freqList[i].ip).Set(float64(freqList[i].count))
	}
}

func isPrivateIP(ip string) bool {
	privateBlocks := []string{"10.", "192.168.", "172.16.", "172.17.", "172.18.", "172.19.", "172.20.", "172.21.",
		"172.22.", "172.23.", "172.24.", "172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
		"172.30.", "172.31."}
	for _, prefix := range privateBlocks {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	return false
}

