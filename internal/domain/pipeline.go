package domain

// Pipeline - доменная сущность пайплайна команд
type Pipeline struct {
	Commands []*Command
	Operator string
}

// NewPipeline создает новый пайплайн
func NewPipeline() *Pipeline {
	return &Pipeline{
		Commands: make([]*Command, 0),
	}
}

// AddCommand добавляет команду в пайплайн
func (p *Pipeline) AddCommand(cmd *Command) {
	p.Commands = append(p.Commands, cmd)
}

// SetOperator устанавливает оператор между командами
func (p *Pipeline) SetOperator(op string) {
	p.Operator = op
}

// IsSingleCommand проверяет, является ли пайплайн одиночной командой
func (p *Pipeline) IsSingleCommand() bool {
	return len(p.Commands) == 1
}

// HasOperator проверяет наличие оператора
func (p *Pipeline) HasOperator() bool {
	return p.Operator != ""
}

// ShouldContinueExecution определяет, должно ли выполнение продолжаться
func (p *Pipeline) ShouldContinueExecution(lastExitCode int) bool {
	if !p.HasOperator() {
		return true
	}

	switch p.Operator {
	case "&&":
		return lastExitCode == 0
	case "||":
		return lastExitCode != 0
	default:
		return true
	}
}
