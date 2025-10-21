package output_adapters

import (
	"minishell/internal/application/ports"
	"minishell/internal/domain"
)

// CommandExecutorAdapter - выходной адаптер для выполнения команд
type CommandExecutorAdapter struct {
	commandService ports.CommandInputPort
}

// NewCommandExecutorAdapter создает новый адаптер исполнителя
func NewCommandExecutorAdapter(
	commandService ports.CommandInputPort,
) *CommandExecutorAdapter {
	return &CommandExecutorAdapter{
		commandService: commandService,
	}
}

// ExecutePipeline выполняет пайплайн команд
func (e *CommandExecutorAdapter) ExecutePipeline(pipeline *domain.Pipeline, ctx *domain.ExecutionContext) error {
	return e.commandService.ExecutePipeline(pipeline, ctx)
}

// ExecuteSingleCommand выполняет одиночную команду
func (e *CommandExecutorAdapter) ExecuteSingleCommand(cmd *domain.Command, ctx *domain.ExecutionContext) error {
	return e.commandService.ExecuteSingleCommand(cmd, ctx)
}
