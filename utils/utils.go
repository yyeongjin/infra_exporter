package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RunCommand는 시스템 명령어를 실행하고, 출력과 오류를 반환합니다.
func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("명령어 실행 실패: %v, 출력: %s", err, out.String())
	}
	return out.String(), nil
}

// GetCurrentTime는 현재 시간을 RFC3339 형식으로 반환합니다.
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// FileExists는 주어진 경로에 파일이 존재하는지 확인합니다. 
// 파일이 존재하지 않거나 다른 오류가 발생하면 false를 반환합니다.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	// 파일이 존재하지 않으면 false를 반환하고, 다른 오류는 로깅
	if os.IsNotExist(err) {
		return false
	}
	// 다른 오류가 발생하면 로그를 기록하고 false 반환
	fmt.Printf("파일 존재 여부 확인 중 오류: %v\n", err)
	return false
}

// CreateLogFile은 로그 파일을 생성하고 로그를 기록합니다.
// 로그 파일을 열고, 로그 메시지를 기록합니다.
func CreateLogFile(filename, message string) error {
	file, err := openLogFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 로그 작성
	_, err = file.WriteString(fmt.Sprintf("[%s] %s\n", GetCurrentTime(), message))
	if err != nil {
		return fmt.Errorf("로그 파일 쓰기 실패: %v", err)
	}
	return nil
}

// openLogFile은 로그 파일을 열거나 새로 생성하는 유틸리티 함수입니다.
func openLogFile(filename string) (*os.File, error) {
	// 파일을 열거나 새로 생성
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("로그 파일 열기 실패: %v", err)
	}
	return file, nil
}

// GetUserID는 사용자 계정의 ID를 반환합니다. 공백을 제거한 ID 값을 반환합니다.
func GetUserID(username string) (string, error) {
	cmdOutput, err := RunCommand("id", "-u", username)
	if err != nil {
		return "", fmt.Errorf("사용자 ID를 얻는 데 실패: %v", err)
	}
	// 출력에서 앞뒤 공백을 제거하여 반환
	return strings.TrimSpace(cmdOutput), nil
}

// GetUserGroupID는 사용자 계정의 그룹 ID를 반환합니다. 공백을 제거한 그룹 ID 값을 반환합니다.
func GetUserGroupID(username string) (string, error) {
	cmdOutput, err := RunCommand("id", "-g", username)
	if err != nil {
		return "", fmt.Errorf("사용자 그룹 ID를 얻는 데 실패: %v", err)
	}
	// 출력에서 앞뒤 공백을 제거하여 반환
	return strings.TrimSpace(cmdOutput), nil
}

