### Minishell: взаимодействие с ОС
Необходимо реализовать собственный простейший Unix shell.

Требования
Ваш интерпретатор командной строки должен поддерживать:

Встроенные команды:
– cd <path> – смена текущей директории.
– pwd – вывод текущей директории.
– echo <args> – вывод аргументов.
– kill <pid> – послать сигнал завершения процессу с заданным PID.
– ps – вывести список запущенных процессов.

Запуск внешних команд через exec (с помощью системных вызовов fork/exec либо стандартных функций os/exec).

Конвейеры (pipelines): возможность объединять команды через |, чтобы вывод одной команды направлять на ввод следующей (как в обычном shell).

Например: ps | grep myprocess | wc -l.

Обработку завершения: при нажатии Ctrl+D (EOF) шелл должен завершаться; Ctrl+C — прерывание текущей запущенной команды, но без закрыватия самой shell.

Дополнительно: реализовать парсинг && и || (условное выполнение команд), подстановку переменных окружения $VAR, поддержку редиректов >/< для вывода в файл и чтения из файла.

Основной упор необходимо делать на реализацию базового функционала (exec, builtins, pipelines). Проверять надо как интерактивно, так и скриптом. Код должен работать без ситуаций гонки, корректно освобождать ресурсы. 

Совет: используйте пакеты os/exec, bufio (для ввода), strings.Fields (для разбиения командной строки на аргументы) и системные вызовы через syscall, если потребуется.

Структура проекта:
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
│       ├── adapters/
│       │   ├── input_adapters/
│       │   │   └── shell_controller.go
│       │   ├── output_adapters/
│       │   │   ├── command_executor_adapter.go
│       │   │   └── system_repository_adapter.go
│       │   ├── parser_adapters/
│       │   │   └── command_parser_adapter.go
│       │   └── presenters/
│       │       └── shell_presenter_adapter.go
│       └── frameworks/
├── pkg/
│   ├── utils/
│   │   └── string_utils.go
│   └── constants/
│       └── shell_constants.go
└── go.mod
```