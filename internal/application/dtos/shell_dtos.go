package dtos

// ShellConfigDTO - DTO для конфигурации shell
type ShellConfigDTO struct {
	Prompt    string
	IsRunning bool
}

// ExecutionContextDTO - DTO для контекста выполнения
type ExecutionContextDTO struct {
	CurrentDir   string
	LastExitCode int
	Environment  map[string]string
}
