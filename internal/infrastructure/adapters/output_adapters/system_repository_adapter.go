package output_adapters

import (
	"bytes"
	"fmt"
	"io"
	"minishell/internal/domain"
	"minishell/pkg/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// SystemRepositoryAdapter - выходной адаптер для системных операций
type SystemRepositoryAdapter struct{}

// NewSystemRepositoryAdapter создает новый адаптер системного репозитория
func NewSystemRepositoryAdapter() *SystemRepositoryAdapter {
	return &SystemRepositoryAdapter{}
}

// ExecuteCommand выполняет внешнюю команду
func (r *SystemRepositoryAdapter) ExecuteCommand(cmd *domain.Command, input []byte) ([]byte, int, error) {
	execCmd := exec.Command(cmd.Name, cmd.Args...)

	var stdin io.Reader
	var stdout, stderr bytes.Buffer

	if input != nil {
		stdin = bytes.NewReader(input)
	} else if cmd.Input != "" {
		file, err := os.Open(cmd.Input)
		if err != nil {
			return nil, 1, err
		}
		defer file.Close()
		stdin = file
	} else {
		stdin = os.Stdin
	}

	execCmd.Stdin = stdin
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	execCmd.Env = r.getEnvironment()

	if cmd.Output != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if cmd.Append {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		outputFile, err := os.OpenFile(cmd.Output, flags, 0644)
		if err != nil {
			return nil, 1, err
		}
		defer outputFile.Close()

		execCmd.Stdout = outputFile
	}

	var exitCode int
	err := execCmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	} else {
		exitCode = 0
	}

	if stderr.Len() > 0 {
		fmt.Fprint(os.Stderr, stderr.String())
	}

	return stdout.Bytes(), exitCode, nil
}

// ChangeDirectory меняет текущую директорию
func (r *SystemRepositoryAdapter) ChangeDirectory(path string) error {
	return os.Chdir(path)
}

// GetCurrentDirectory возвращает текущую директорию
func (r *SystemRepositoryAdapter) GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

// GetEnvironment возвращает переменные окружения
func (r *SystemRepositoryAdapter) GetEnvironment() map[string]string {
	env := make(map[string]string)
	for _, pair := range os.Environ() {
		parts := utils.SplitEnvVar(pair)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

// getEnvironment возвращает окружение в формате для exec
func (r *SystemRepositoryAdapter) getEnvironment() []string {
	return os.Environ()
}

// KillProcess убивает процесс по PID
func (r *SystemRepositoryAdapter) KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGTERM)
}

// GetProcessList возвращает список процессов
func (r *SystemRepositoryAdapter) GetProcessList() ([]domain.ProcessInfo, error) {
	var processes []domain.ProcessInfo

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if pid, err := strconv.Atoi(entry.Name()); err == nil {
			cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
			if cmdline, err := os.ReadFile(cmdlinePath); err == nil {
				cmdStr := strings.ReplaceAll(string(cmdline), "\x00", " ")
				cmdStr = strings.TrimSpace(cmdStr)
				if cmdStr != "" {
					processes = append(processes, domain.ProcessInfo{
						PID: pid,
						Cmd: cmdStr,
					})
				}
			}
		}
	}

	return processes, nil
}
