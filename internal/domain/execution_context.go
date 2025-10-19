package domain

// NewExecutionContext создает новый контекст выполнения
func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		Environment:  make(map[string]string),
		IsRunning:    true,
		LastExitCode: 0,
	}
}

// UpdateCurrentDir обновляет текущую директорию
func (ctx *ExecutionContext) UpdateCurrentDir(dir string) {
	ctx.CurrentDir = dir
}

// SetEnv устанавливает переменную окружения
func (ctx *ExecutionContext) SetEnv(key, value string) {
	ctx.Environment[key] = value
}

// GetEnv получает переменную окружения
func (ctx *ExecutionContext) GetEnv(key string) string {
	return ctx.Environment[key]
}

// UpdateExitCode обновляет код завершения последней команды
func (ctx *ExecutionContext) UpdateExitCode(code int) {
	ctx.LastExitCode = code
}

// Stop останавливает выполнение shell
func (ctx *ExecutionContext) Stop() {
	ctx.IsRunning = false
}

// GetPrompt возвращает строку приглашения
func (ctx *ExecutionContext) GetPrompt() string {
	return "minishell:" + ctx.CurrentDir + "$ "
}
