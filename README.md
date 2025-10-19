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
│   ├── domain/             # Domain Layer - бизнес-сущности и правила
│   │   ├── command.go
|   |   ├── entities.go
│   │   ├── pipeline.go
│   │   └── context.go
│   ├── application/        # Application Layer - use cases и бизнес-логика
│   │   ├── ports/
│   │   │   ├── input.go
│   │   │   └── output.go
│   │   ├── services/
│   │   │   ├── shell_service.go
│   │   │   └── command_service.go
│   │   └── dtos/
│   │       └── commands.go
│   ├── interfaces/         # Interface Adapters Layer
│   │   ├── controllers/
│   │   │   └── shell_controller.go
│   │   ├── repositories/
│   │   │   └── system_repository.go
│   │   └── presenters/
│   │       └── shell_presenter.go
│   └── infrastructure/     # Infrastructure Layer (Frameworks & Drivers)
│       ├── parser/
│       │   └── command_parser.go
│       ├── executors/
│       │   └── command_executor.go
│       └── system/
│           └── os_operations.go
└── go.mod
```