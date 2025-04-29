package collector

import (
    "log"
    "os"
    "sync"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

var (
    sensitiveFileMetric = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "infra_sensitive_file_changes_total",
            Help: "민감 파일 변경 감지 횟수",
        },
        []string{"file"},
    )

    sensitiveFileMutex sync.Mutex
    lastModifiedMap     = make(map[string]time.Time)
)

func RegisterSensitiveFileCollector(paths []string) {
    prometheus.MustRegister(sensitiveFileMetric)

    for _, path := range paths {
        go func(p string) {
            for {
                checkSensitiveFileChange(p)
                time.Sleep(1 * time.Minute) // 1분마다 파일 변경 감지
            }
        }(path)
    }
}

func checkSensitiveFileChange(file string) {
    sensitiveFileMutex.Lock()
    defer sensitiveFileMutex.Unlock()

    info, err := os.Stat(file)
    if err != nil {
        log.Printf("파일 정보 가져오기 실패 (%s): %v", file, err)
        return
    }

    modTime := info.ModTime().Local() // 파일 수정 시간 (로컬 시간대)
    currentTime := time.Now().Local()  // 현재 시간 (로컬 시간대)

    lastTime, ok := lastModifiedMap[file]
    if !ok {
        // 파일을 처음 읽은 경우
        lastModifiedMap[file] = modTime
        log.Printf("파일 %s 최초 수정 시간 저장: %s", file, modTime.Format("2006-01-02 15:04:05.000000000 -0700 MST"))
        return
    }

    // 파일이 수정되었는지 확인
    if !modTime.Equal(lastTime) {
        lastModifiedMap[file] = modTime
        sensitiveFileMetric.WithLabelValues(file).Inc()
        log.Printf("민감 파일 변경 감지: %s (변경 시간: %s)", file, modTime.Format("2006-01-02 15:04:05.000000000 -0700 MST"))
    } else {
        log.Printf("파일 %s는 변경되지 않았습니다 (현재 시간: %s)", file, currentTime.Format("2006-01-02 15:04:05.000000000 -0700 MST"))
    }
}

