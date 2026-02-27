package commands

import "fmt"

// Command define a interface que todos os subcomandos devem implementar
type Command interface {
	Name() string
	Aliases() []string
	Description() string
	Run(args []string) error
}

// Registry armazena todos os comandos disponíveis (por nome e por alias)
var Registry = make(map[string]Command)

// CommandsList mantém apenas os comandos registrados (sem aliases) para exibição na ajuda
var CommandsList []Command

func Register(cmd Command) {
	// Registra o nome principal
	Registry[cmd.Name()] = cmd

	// Registra todos os aliases apontando para a mesma estrutura
	for _, alias := range cmd.Aliases() {
		if _, exists := Registry[alias]; exists {
			panic(fmt.Sprintf("Erro de colisão: O alias/nome '%s' já está em uso!", alias))
		}
		Registry[alias] = cmd
	}

	// Guarda na lista para o menu de ajuda
	CommandsList = append(CommandsList, cmd)
}
