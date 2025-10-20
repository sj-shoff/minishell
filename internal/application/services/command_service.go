package services

import (
	"bytes"
	"fmt"
	"minishell/internal/application/ports"
	"minishell/internal/domain"
	"strconv"
)

// CommandService - application service для выполнения команд
type CommandService struct {
	system ports.SystemRepositoryOutputPort
}

// NewCommandService создает новый сервис команд
func NewCommandService(system ports.SystemRepositoryOutputPort) *CommandService {
	return &CommandService{
		system: system,
	}
}

// ExecutePipeline выполняет пайплайн команд
func (s *CommandService) ExecutePipeline(pipeline *domain.Pipeline, ctx *domain.ExecutionContext) error {
	if pipeline.IsSingleCommand() {
		s.ExecuteSingleCommand(pipeline.Commands[0], ctx)
	}
	return s.executePipeSequence(pipeline.Commands, ctx)
}

// ExecuteSingleCommand выполняет одиночную команду
func (s *CommandService) ExecuteSingleCommand(cmd *domain.Command, ctx *domain.ExecutionContext) error {
	if cmd.IsBuiltin() {
		return s.executeBuiltinCommand(cmd, ctx)
	}

	output, exitCode, err := s.system.ExecuteCommand(cmd, nil)
	if err != nil {
		return err
	}

	ctx.UpdateExitCode(exitCode)

	if len(output) > 0 {
		fmt.Print(string(output))
	}

	return nil
}

// executePipeSequence выполняет последовательность команд с пайпами
func (s *CommandService) executePipeSequence(commands []*domain.Command, ctx *domain.ExecutionContext) error {
	var input []byte
	var lastExitCode int

	for i, cmd := range commands {
		output, exitcode, err := s.system.ExecuteCommand(cmd, input)
		if err != nil {
			ctx.UpdateExitCode(1)
		}

		input = output
		lastExitCode = exitcode

		if i == len(commands)-1 && len(output) > 0 {
			fmt.Print(string(output))
		}
	}

	ctx.UpdateExitCode(lastExitCode)
	return nil
}

// executeBuiltinCommand выполняет встроенную команду
func (s *CommandService) executeBuiltinCommand(cmd *domain.Command, ctx *domain.ExecutionContext) error {
	switch cmd.Name {
	case "cd":
		return s.executeCD(cmd, ctx)
	case "pwd":
		return s.executePWD(ctx)
	case "echo":
		return s.executeEcho(cmd, ctx)
	case "kill":
		return s.executeKill(cmd, ctx)
	case "ps":
		return s.executePS(ctx)
	default:
		ctx.UpdateExitCode(1)
		return fmt.Errorf("unknown builtin command: %s", cmd.Name)
	}
}

// executeCD выполняет команду cd
func (s *CommandService) executeCD(cmd *domain.Command, ctx *domain.ExecutionContext) error {
	var path string
	if len(cmd.Args) == 0 {
		path = ctx.GetEnv("HOME")
		if path == "" {
			ctx.UpdateExitCode(1)
			return fmt.Errorf("HOME not set")
		}
	} else {
		path = cmd.Args[0]
	}

	if err := s.system.ChangeDirectory(path); err != nil {
		ctx.UpdateExitCode(1)
		return err
	}

	if dir, err := s.system.GetCurrentDirectory(); err == nil {
		ctx.UpdateCurrentDir(dir)
	}

	return nil
}

// executePWD выполняет команду pwd
func (s *CommandService) executePWD(ctx *domain.ExecutionContext) error {
	dir, err := s.system.GetCurrentDirectory()
	if err != nil {
		ctx.UpdateExitCode(1)
		return err
	}

	fmt.Print(dir)
	ctx.UpdateExitCode(0)
	return nil
}

// executeEcho выполняет команду echo
func (s *CommandService) executeEcho(cmd *domain.Command, ctx *domain.ExecutionContext) error {
	output := bytes.Buffer{}
	for i, arg := range cmd.Args {
		if i > 0 {
			output.WriteString(" ")
		}
		output.WriteString(arg)
	}

	fmt.Print(output.String())
	ctx.UpdateExitCode(0)
	return nil
}

// executeKill выполняет команду kill
func (s *CommandService) executeKill(cmd *domain.Command, ctx *domain.ExecutionContext) error {
	if len(cmd.Args) == 0 {
		ctx.UpdateExitCode(1)
		return fmt.Errorf("kill: missing pid")
	}

	pid, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		ctx.UpdateExitCode(1)
		return fmt.Errorf("kill: invalid pid")
	}
	if err := s.system.KillProcess(pid); err != nil {
		ctx.UpdateExitCode(1)
		return err
	}

	ctx.UpdateExitCode(0)
	return nil
}

// executePS выполняет команду ps
func (s *CommandService) executePS(ctx *domain.ExecutionContext) error {
	processes, err := s.system.GetProcessList()
	if err != nil {
		ctx.UpdateExitCode(1)
		return err
	}

	fmt.Println("PID\tCMD")
	for _, proc := range processes {
		fmt.Printf("%d\t%s\n", proc.PID, proc.Cmd)
	}

	ctx.UpdateExitCode(0)
	return nil
}
