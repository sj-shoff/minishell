package presenters

import (
	"fmt"
	"os"
)

// ShellPresenterAdapter - адаптер для представления результатов
type ShellPresenterAdapter struct{}

// NewShellPresenterAdapter создает новый адаптер презентера
func NewShellPresenterAdapter() *ShellPresenterAdapter {
	return &ShellPresenterAdapter{}
}

// ShowPrompt показывает приглашение командной строки
func (p *ShellPresenterAdapter) ShowPrompt(prompt string) {
	fmt.Print(prompt)
}

// ShowOutput показывает вывод команды
func (p *ShellPresenterAdapter) ShowOutput(output string) {
	fmt.Println(output)
}

// ShowError показывает ошибку
func (p *ShellPresenterAdapter) ShowError(error string) {
	fmt.Fprintln(os.Stderr, "Error:", error)
}

// ShowExitCode показывает код завершения
func (p *ShellPresenterAdapter) ShowExitCode(code int) {
	if code != 0 {
		fmt.Fprintf(os.Stderr, "Command exited with code: %d\n", code)
	}
}
