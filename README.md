### Minishell: Реализация Unix Shell

Minishell - это собственная реализация Unix shell, написанная на Go. Проект разработан с использованием Domain-Driven Design (DDD) архитектуры и реализует полностью функциональный интерпретатор командной строки, который поддерживает встроенные команды, выполнение внешних команд, конвейеры, логические операторы, переменные окружения и перенаправления ввода/вывода.

## Поддерживаемые возможности

Встроенные команды:

- cd <path> - смена текущей директории
- pwd - вывод текущей директории
- echo <args> - вывод аргументов
- kill <pid> - послать сигнал завершения процессу с заданным PID
- ps - вывести список запущенных процессов

Внешние команды:

- Выполнение через пакет os/exec
- Поддержка всех системных команд, доступных в PATH

Конвейеры (pipelines):

- Объединение команд с помощью оператора |
- Перенаправление вывода между командами
- Пример: ps | grep myprocess | wc -l

Логические операторы:

- && - условное И (выполнение следующей команды только при успешном завершении предыдущей)
- || - условное ИЛИ (выполнение следующей команды только при неуспешном завершении предыдущей)

Переменные окружения:

Подстановка $VAR в командах
Доступ к переменным окружения системы
Перенаправления ввода/вывода:
```
> - вывод в файл (перезапись)
>> - вывод в файл (добавление)
< - ввод из файла
```
Обработка сигналов:

- Ctrl+D (EOF) - завершение shell
- Ctrl+C - прерывание текущей команды без выхода из shell


## Требования
- Go 1.24+
- Операционная система Unix (Linux, macOS)

## Запуск и тестирование проекта
```bash
# Запуск проекта
make run

# Тестирование скриптом
make sh_test
```

## Структура проекта:
```
minishell/
├── cmd/
│   └── minishell/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── command.go
|   |   ├── execution_context.go
│   │   ├── pipeline.go
│   │   └── process.go
│   ├── application/
│   │   ├── ports/
│   │   │   ├── input_ports.go
│   │   │   └── output_ports.go
│   │   ├── services/
│   │   │   ├── shell_service.go
│   │   │   └── command_service.go
│   │   └── dtos/
│   │       ├── command_dtos.go
│   │       └── shell_dtos.go
│   └── infrastructure/
│       └── adapters/
│           ├── input_adapters/
│           │   └── shell_controller.go
│           ├── output_adapters/
│           │   ├── command_executor_adapter.go
│           │   └── system_repository_adapter.go
│           ├── parser_adapters/
│           │   └── command_parser_adapter.go
│           └── presenters/
│               └── shell_presenter_adapter.go
├── pkg/
│   ├── utils/
│   │   └── string_utils.go
│   └── constants/
│       └── shell_constants.go
├── go.mod
├── Makefile
└── README.md

```
