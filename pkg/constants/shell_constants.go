package constants

const (
	// Exit codes
	ExitSuccess = 0
	ExitFailure = 1

	// Operators
	OperatorPipe   = "|"
	OperatorAnd    = "&&"
	OperatorOr     = "||"
	OperatorOutput = ">"
	OperatorAppend = ">>"
	OperatorInput  = "<"

	// Builtin commands
	CmdCD   = "cd"
	CmdPWD  = "pwd"
	CmdEcho = "echo"
	CmdKill = "kill"
	CmdPS   = "ps"
	CmdExit = "exit"
)
