package input_adapters

import (
	"bufio"
	"fmt"
	"io"
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
	context *domain.ExecutionContext,
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

	reader := bufio.NewReader(os.Stdin)

	for c.shellService.ShouldContinue(c.context) {
		fmt.Print(c.shellService.GetPrompt(c.context))

		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nExit")
				break
			}
			fmt.Println("Read error:", err.Error())
			continue
		}

		input = input[:len(input)-1]
		if input == "" {
			continue
		}

		if err := c.shellService.ExecuteCommand(input, c.context); err != nil {
			fmt.Println("Error:", err.Error())
		}
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
				os.Exit(0)
			case syscall.SIGTERM:
				c.shellService.ExecuteCommand("exit", c.context)
			}
		}
	}()
}
