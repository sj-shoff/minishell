package parser_adapters

import (
	"minishell/internal/domain"
	"minishell/pkg/utils"
	"os"
	"strings"
)

// CommandParserAdapter - адаптер для парсинга команд
type CommandParserAdapter struct{}

// NewCommandParserAdapter создает новый адаптер парсера
func NewCommandParserAdapter() *CommandParserAdapter {
	return &CommandParserAdapter{}
}

// Parse разбирает строку команды на пайплайны
func (p *CommandParserAdapter) Parse(input string, env map[string]string) ([]*domain.Pipeline, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	var pipelines []*domain.Pipeline

	parts := p.splitByLogicalOperators(input)

	// Первый пайплайн всегда выполняется
	if len(parts) > 0 {
		pipeline, err := p.parsePipeline(parts[0], env)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, pipeline)
	}

	// Обрабатываем остальные пайплайны с операторами
	for i := 1; i < len(parts); i += 2 {
		if i+1 >= len(parts) {
			return nil, &ParseError{"missing command after operator"}
		}

		operator := parts[i]
		pipelineStr := parts[i+1]

		pipeline, err := p.parsePipeline(pipelineStr, env)
		if err != nil {
			return nil, err
		}

		pipeline.SetOperator(operator)
		pipelines = append(pipelines, pipeline)
	}

	return pipelines, nil
}

// splitByLogicalOperators разделяет строку по && и ||
func (p *CommandParserAdapter) splitByLogicalOperators(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(' ')

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch == '"' || ch == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuotes = false
			}
		}

		if !inQuotes && i < len(input)-1 {
			op := input[i : i+2]
			if op == "&&" || op == "||" {
				if current.Len() > 0 {
					parts = append(parts, utils.TrimSpace(current.String()))
					parts = append(parts, op)
					current.Reset()
					i++
					continue
				}
			}
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		parts = append(parts, utils.TrimSpace(current.String()))
	}

	return parts
}

// parsePipeline разбирает пайплайн команд
func (p *CommandParserAdapter) parsePipeline(pipelineStr string, env map[string]string) (*domain.Pipeline, error) {
	commands := utils.SplitFields(pipelineStr, "|")
	pipeline := domain.NewPipeline()

	for _, cmdStr := range commands {
		cmd, err := p.parseCommand(utils.TrimSpace(cmdStr), env)
		if err != nil {
			return nil, err
		}
		pipeline.AddCommand(cmd)
	}

	return pipeline, nil
}

// parseCommand разбирает одну команду
func (p *CommandParserAdapter) parseCommand(cmdStr string, env map[string]string) (*domain.Command, error) {
	tokens := p.tokenize(cmdStr)
	if len(tokens) == 0 {
		return nil, nil
	}

	cmd := domain.NewCommand(p.expandVariables(tokens[0], env))
	i := 1

	for i < len(tokens) {
		token := p.expandVariables(tokens[i], env)

		switch token {
		case ">", ">>", "<":
			if i+1 >= len(tokens) {
				return nil, &ParseError{"missing filename for redirection"}
			}

			filename := p.expandVariables(tokens[i+1], env)
			switch token {
			case ">":
				cmd.SetOutput(filename, false)
			case ">>":
				cmd.SetOutput(filename, true)
			case "<":
				cmd.SetInput(filename)
			}
			i += 2   // Пропускаем и оператор и имя файла
			continue // Не добавляем оператор как аргумент

		case "&":
			cmd.Background = true
			i++
			continue

		default:
			cmd.AddArg(token)
			i++
		}
	}

	return cmd, nil
}

// tokenize разбивает строку на токены с учетом кавычек
func (p *CommandParserAdapter) tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(' ')

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch {
		case ch == '"' || ch == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuotes = false
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(ch)
			}

		case ch == ' ' && !inQuotes:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// expandVariables заменяет переменные окружения в строке
func (p *CommandParserAdapter) expandVariables(input string, env map[string]string) string {
	return os.Expand(input, func(key string) string {
		if value, exists := env[key]; exists {
			return value
		}
		return os.Getenv(key)
	})
}

// ParseError представляет ошибку парсинга
type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}
