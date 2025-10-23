package output_adapters

import (
	"bytes"
	"fmt"
	"io"
	"minishell/internal/domain"
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
	// Проверяем, существует ли команда
	if cmd.Name == "" {
		return nil, 127, fmt.Errorf("command not found")
	}

	// Проверяем, доступна ли команда в системе
	if _, err := exec.LookPath(cmd.Name); err != nil {
		return nil, 127, fmt.Errorf("command not found: %s", cmd.Name)
	}

	execCmd := exec.Command(cmd.Name, cmd.Args...)

	var stdin io.Reader
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Настройка ввода
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

	// Настройка вывода - по умолчанию в буфер
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	// Обработка редиректа вывода
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

		// Перенаправляем вывод в файл
		execCmd.Stdout = outputFile
		// Для редиректа вывода не возвращаем данные в stdout
		stdout.Reset()
	}

	execCmd.Env = os.Environ()

	var exitCode int
	err := execCmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// Если команда не найдена, возвращаем код 127
			exitCode = 127
		}
	} else {
		exitCode = 0
	}

	// Выводим stderr только если нет редиректа вывода
	if stderr.Len() > 0 && cmd.Output == "" {
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
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
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

	// Используем системную команду ps для получения списка процессов
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || line == "" { // Пропускаем заголовок и пустые строки
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if pid, err := strconv.Atoi(fields[1]); err == nil {
			processes = append(processes, domain.ProcessInfo{
				PID: pid,
				Cmd: strings.Join(fields[10:], " "),
			})
		}
	}

	return processes, nil
}
