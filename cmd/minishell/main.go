package main

import (
	"minishell/internal/application/services"
	input_adapters "minishell/internal/infrastructure/adapters/input_adapters"
	output_adapters "minishell/internal/infrastructure/adapters/output_adapters"
	parser_adapters "minishell/internal/infrastructure/adapters/parser_adapters"
	presenters "minishell/internal/infrastructure/adapters/presenters"
)

func main() {
	// Инициализация адаптеров
	systemRepo := output_adapters.NewSystemRepositoryAdapter()
	commandParser := parser_adapters.NewCommandParserAdapter()
	shellPresenter := presenters.NewShellPresenterAdapter()

	// Инициализация сервисов
	commandService := services.NewCommandService(systemRepo)
	shellService := services.NewShellService(
		commandParser,
		commandService,
		systemRepo,
		shellPresenter,
	)

	// Инициализация контроллера
	shellController := input_adapters.NewShellController(shellService, systemRepo)

	// Запуск приложения
	shellController.Run()
}
