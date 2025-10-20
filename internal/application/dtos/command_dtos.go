package dtos

// CommandDTO - DTO для передачи данных команды
type CommandDTO struct {
	Name       string
	Args       []string
	Input      string
	Output     string
	Append     bool
	Background bool
}

// PipelineDTO - DTO для передачи данных пайплайна
type PipelineDTO struct {
	Commands []CommandDTO
	Operator string
}

// CommandResultDTO - DTO для передачи результата команды
type CommandResultDTO struct {
	ExitCode int
	Output   string
	Error    string
}
