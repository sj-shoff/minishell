package ports

import "minishell/internal/domain"

// CommandParserOutputPort - исходящий порт для парсинга команд
type CommandParserOutputPort interface {
	Parse(input string, env map[string]string) ([]*domain.Pipeline, error)
}

// SystemRepositoryOutputPort - исходящий порт для системных операций
type SystemRepositoryOutputPort interface {
	ExecuteCommand(cmd *domain.Command, input []byte) ([]byte, int, error)
	ChangeDirectory(path string) error
	GetCurrentDirectory() (string, error)
	GetEnvironment() map[string]string
	KillProcess(pid int) error
	GetProcessList() ([]domain.ProcessInfo, error)
}

// ShellPresenterOutputPort - исходящий порт для представления результатов
type ShellPresenterOutputPort interface {
	ShowPrompt(prompt string)
	ShowOutput(output string)
	ShowError(error string)
	ShowExitCode(code int)
}
