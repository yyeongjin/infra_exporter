package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// Config는 애플리케이션의 구성 설정을 나타내는 구조체입니다.
type Config struct {
	Monitor MonitorConfig `yaml:"monitor"`
}

// MonitorConfig는 모니터링 설정을 나타냅니다.
type MonitorConfig struct {
	SSHFailed     bool               `yaml:"ssh_failed"`
	SudoUsage     bool               `yaml:"sudo_usage"`
	UserStatus    bool               `yaml:"user_status"`
	Ports         bool               `yaml:"ports"`
	ExternalIP    ExternalIPConfig   `yaml:"external_ip"`
	SensitiveFile SensitiveFileConf  `yaml:"sensitive_file"`
}

// ExternalIPConfig는 외부 IP 수집 설정을 나타냅니다.
type ExternalIPConfig struct {
	Enabled bool `yaml:"enabled"`
	TopN    int  `yaml:"top_n"`
}

// SensitiveFileConf는 민감한 파일 수집 설정을 나타냅니다.
type SensitiveFileConf struct {
	Enabled bool     `yaml:"enabled"`
	Paths   []string `yaml:"paths"`
}

// LoadConfig는 지정된 경로에서 YAML 구성을 로드합니다.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("구성 파일 읽기 실패: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("구성 파일 파싱 실패: %w", err)
	}

	// 유효성 검사
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	// 환경 변수로 오버라이드 (필요시)
	overrideConfigWithEnv(&cfg)

	return &cfg, nil
}

// validateConfig는 구성 설정의 유효성을 검사합니다.
func validateConfig(cfg *Config) error {
	// SSH 실패 추적 여부가 설정되지 않으면 기본값을 true로 설정
	if !cfg.Monitor.SSHFailed && !cfg.Monitor.SudoUsage && !cfg.Monitor.UserStatus && !cfg.Monitor.Ports && !cfg.Monitor.ExternalIP.Enabled && !cfg.Monitor.SensitiveFile.Enabled {
		return fmt.Errorf("적어도 하나의 모니터링 항목이 활성화되어야 합니다")
	}
	
	// ExternalIP 설정에서 TopN이 0일 경우 기본값을 5로 설정
	if cfg.Monitor.ExternalIP.Enabled && cfg.Monitor.ExternalIP.TopN <= 0 {
		cfg.Monitor.ExternalIP.TopN = 5
	}
	return nil
}

// overrideConfigWithEnv는 환경 변수로 구성을 오버라이드합니다.
func overrideConfigWithEnv(cfg *Config) {
	// 환경 변수로 설정 값을 오버라이드 (예: MONITOR_SSH_FAILED)
	if value := os.Getenv("MONITOR_SSH_FAILED"); value != "" {
		cfg.Monitor.SSHFailed = parseBool(value)
	}
	if value := os.Getenv("MONITOR_SUDO_USAGE"); value != "" {
		cfg.Monitor.SudoUsage = parseBool(value)
	}
	if value := os.Getenv("MONITOR_USER_STATUS"); value != "" {
		cfg.Monitor.UserStatus = parseBool(value)
	}
	if value := os.Getenv("MONITOR_PORTS"); value != "" {
		cfg.Monitor.Ports = parseBool(value)
	}
	if value := os.Getenv("MONITOR_EXTERNAL_IP_ENABLED"); value != "" {
		cfg.Monitor.ExternalIP.Enabled = parseBool(value)
	}
	if value := os.Getenv("MONITOR_SENSITIVE_FILE_ENABLED"); value != "" {
		cfg.Monitor.SensitiveFile.Enabled = parseBool(value)
	}
}

// parseBool은 문자열을 불리언 값으로 변환합니다.
func parseBool(value string) bool {
	value = strings.ToLower(value)
	return value == "true" || value == "1"
}

