package ports

import "minishell/internal/domain"

// ShellInputPort - входящий порт для операций shell
type ShellInputPort interface {
	ExecuteCommand(input string, ctx *domain.ExecutionContext) error
	ShouldContinue(ctx *domain.ExecutionContext) bool
	GetPrompt(ctx *domain.ExecutionContext) string
}

// CommandInputPort - входящий порт для выполнения команд
type CommandInputPort interface {
	ExecutePipeline(pipeline *domain.Pipeline, ctx *domain.ExecutionContext) error
	ExecuteSingleCommand(cmd *domain.Command, ctx *domain.ExecutionContext) error
}
