package domain

// Command - доменная сущность команды
type Command struct {
	Name       string
	Args       []string
	Input      string
	Output     string
	Append     bool
	Background bool
}

// NewCommand создает новую команду
func NewCommand(name string) *Command {
	return &Command{
		Name: name,
		Args: make([]string, 0),
	}
}

// AddArg добавляет аргумент к команде
func (c *Command) AddArg(arg string) {
	c.Args = append(c.Args, arg)
}

// SetInput устанавливает перенаправление ввода
func (c *Command) SetInput(file string) {
	c.Input = file
}

// SetOutput устанавливает перенаправление вывода
func (c *Command) SetOutput(file string, append bool) {
	c.Output = file
	c.Append = append
}

// IsBuiltin проверяет, является ли команда встроенной
func (c *Command) IsBuiltin() bool {
	builtins := map[string]bool{
		"cd":   true,
		"pwd":  true,
		"echo": true,
		"kill": true,
		"ps":   true,
		"exit": true,
	}
	return builtins[c.Name]
}
