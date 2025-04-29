package main

import (
    "log"
    "net/http"
    "os"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "infra_exporter/config"
    "infra_exporter/collector"
)

func main() {
    // config 로딩
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        log.Fatalf("config 로드 실패: %v", err)
    }

    // 기본 메트릭 비활성화
    disableDefaultMetrics()

    // Collector 등록
    registerCollectors(cfg)

    // 메트릭 엔드포인트 오픈
    startMetricsServer()
}

func disableDefaultMetrics() {
    prometheus.Unregister(prometheus.NewGoCollector()) // 기본 Go 메트릭 비활성화
    prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})) // 기본 프로세스 메트릭 비활성화
}

func registerCollectors(cfg *config.Config) {
    // 각 모니터링 항목에 대해 Collector를 등록
    if cfg.Monitor.SSHFailed {
        log.Println("SSH 실패 감지 활성화")
        go collector.RegisterSSHFailedCollector()
    }
    if cfg.Monitor.SudoUsage {
        log.Println("Sudo 사용 감지 활성화")
        go collector.RegisterSudoUsageCollector()
    }
    if cfg.Monitor.SensitiveFile.Enabled {
        log.Println("민감 파일 변경 감지 활성화")
        go collector.RegisterSensitiveFileCollector(cfg.Monitor.SensitiveFile.Paths) // 파일 경로 리스트를 전달
    }
    if cfg.Monitor.ExternalIP.Enabled {
        log.Println("외부 IP 감지 활성화")
        go collector.RegisterExternalIPCollector(cfg.Monitor.ExternalIP.TopN) // TopN 값을 전달
    }
    if cfg.Monitor.UserStatus {
        log.Println("사용자 상태 수집 활성화")
        go collector.RegisterUserStatusCollector()
    }
    if cfg.Monitor.Ports {
        log.Println("리스닝 포트 수집 활성화")
        go collector.RegisterPortCollector()
    }
}

func startMetricsServer() {
    port := getPortFromEnv()
    log.Printf("infra_exporter 실행 중 (:%s/metrics)", port)
    http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPortFromEnv() string {
    port := os.Getenv("INFRA_EXPORTER_PORT")
    if port == "" {
        port = "9101" // 기본 포트 설정
    }
    return port
}

