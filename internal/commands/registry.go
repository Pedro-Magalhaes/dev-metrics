package commands

// Command define a interface que todos os subcomandos devem implementar
type Command interface {
	Name() string
	Description() string
	Run(args []string) error
}

// Registry armazena todos os comandos disponíveis
var Registry = make(map[string]Command)

// Register adiciona um comando ao mapa
func Register(cmd Command) {
	Registry[cmd.Name()] = cmd
}
