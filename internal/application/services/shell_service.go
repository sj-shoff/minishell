package services

import (
	"minishell/internal/application/ports"
	"minishell/internal/domain"
)

// ShellService - application service для операций shell
type ShellService struct {
	parser    ports.CommandParserOutputPort
	executor  ports.CommandInputPort
	system    ports.SystemRepositoryOutputPort
	presenter ports.ShellPresenterOutputPort
}

// NewShellService создает новый сервис shell
func NewShellService(
	parser ports.CommandParserOutputPort,
	executor ports.CommandInputPort,
	system ports.SystemRepositoryOutputPort,
	presenter ports.ShellPresenterOutputPort,
) *ShellService {
	return &ShellService{
		parser:    parser,
		executor:  executor,
		system:    system,
		presenter: presenter,
	}
}

// ExecuteCommand выполняет команду
func (s *ShellService) ExecuteCommand(input string, ctx *domain.ExecutionContext) error {
	if input == "exit" {
		ctx.Stop()
		return nil
	}

	pipelines, err := s.parser.Parse(input, ctx.Environment)
	if err != nil {
		s.presenter.ShowError("Parse error: " + err.Error())
		ctx.UpdateExitCode(1)
		return err
	}

	for _, pipeline := range pipelines {
		if !pipeline.ShouldContinueExecution(ctx.LastExitCode) {
			continue
		}

		err := s.executor.ExecutePipeline(pipeline, ctx)
		if err != nil {
			s.presenter.ShowError("Execution error: " + err.Error())
			break
		}
	}

	if dir, err := s.system.GetCurrentDirectory(); err == nil {
		ctx.UpdateCurrentDir(dir)
	}

	return nil
}

// ShouldContinue проверяет должен ли shell продолжать работу
func (s *ShellService) ShouldContinue(ctx *domain.ExecutionContext) bool {
	return ctx.IsRunning
}

// GetPrompt возвращает строку приглашения
func (s *ShellService) GetPrompt(ctx *domain.ExecutionContext) string {
	return ctx.GetPrompt()
}
