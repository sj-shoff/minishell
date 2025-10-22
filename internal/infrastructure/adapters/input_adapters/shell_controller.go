package input_adapters

import (
	"bufio"
	"fmt"
	"minishell/internal/application/ports"
	"minishell/internal/domain"
	"os"
	"os/signal"
	"syscall"
)

// ShellController - входной адаптер для CLI
type ShellController struct {
	shellService ports.ShellInputPort
	system       ports.SystemRepositoryOutputPort
	context      *domain.ExecutionContext
}

// NewShellController создает новый контроллер
func NewShellController(
	shellService ports.ShellInputPort,
	system ports.SystemRepositoryOutputPort,
) *ShellController {

	ctx := domain.NewExecutionContext()

	if dir, err := system.GetCurrentDirectory(); err == nil {
		ctx.UpdateCurrentDir(dir)
	}

	env := system.GetEnvironment()
	for k, v := range env {
		ctx.SetEnv(k, v)
	}

	return &ShellController{
		shellService: shellService,
		system:       system,
		context:      ctx,
	}
}

// Run запускает основной цикл shell
func (c *ShellController) Run() {
	c.setupSignalHandling()

	scanner := bufio.NewScanner(os.Stdin)

	for c.shellService.ShouldContinue(c.context) {
		fmt.Print(c.shellService.GetPrompt(c.context))

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		if err := c.shellService.ExecuteCommand(input, c.context); err != nil {
			fmt.Println("Error:", err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

// setupSignalHandling настраивает обработку сигналов
func (c *ShellController) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGINT:
				fmt.Println("\nInterrupted")
			case syscall.SIGTERM:
				c.shellService.ExecuteCommand("exit", c.context)
			}
		}
	}()
}
