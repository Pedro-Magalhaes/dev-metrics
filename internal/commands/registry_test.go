package commands_test

import (
	"dev-metrics/internal/commands"
	"fmt"
	"testing"
)

// MockCommand implementa a interface commands.Command
type MockCommand struct {
	name    string
	aliases []string
}

func (m *MockCommand) Name() string            { return m.name }
func (m *MockCommand) Aliases() []string       { return m.aliases }
func (m *MockCommand) Description() string     { return "Descrição de teste" }
func (m *MockCommand) Run(args []string) error { return nil }

func TestRegisterExternal(t *testing.T) {
	// Limpeza do estado global (acessando via pacote)
	commands.Registry = make(map[string]commands.Command)
	commands.CommandsList = []commands.Command{}

	cmd := &MockCommand{
		name:    "help",
		aliases: []string{"h", "ajuda"},
	}

	// Executa a função pública
	commands.Register(cmd)

	// Validação via API pública
	if _, ok := commands.Registry["help"]; !ok {
		t.Error("Falha ao registrar nome principal: 'help' não encontrado no Registry")
	}

	if len(commands.CommandsList) != 1 {
		t.Errorf("Esperava 1 comando na lista, obtive %d", len(commands.CommandsList))
	}
}

func TestPanicOnCollision(t *testing.T) {
	commands.Registry = make(map[string]commands.Command)

	cmd1 := &MockCommand{name: "run", aliases: []string{"r"}}
	cmd2 := &MockCommand{name: "reset", aliases: []string{"r"}} // Conflito aqui

	commands.Register(cmd1)

	// Captura o panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("O código deveria ter causado um panic devido à colisão de aliases")
		} else {
			expectedMsg := "Erro de colisão: O alias/nome 'r' já está em uso!"
			if fmt.Sprint(r) != expectedMsg {
				t.Errorf("Mensagem de erro inesperada: %v", r)
			}
		}
	}()

	commands.Register(cmd2)
}

func TestAliasesPointToSameCommand(t *testing.T) {
	commands.Registry = make(map[string]commands.Command)

	cmd1 := &MockCommand{name: "alpha", aliases: []string{"a"}}
	cmd2 := &MockCommand{name: "run", aliases: []string{"r", "execute"}}
	cmd3 := &MockCommand{name: "zeta", aliases: []string{"z"}}

	cmds := []*MockCommand{cmd1, cmd2, cmd3}
	for _, cmd := range cmds {
		commands.Register(cmd)
	}

	for _, cmd := range cmds {
		for _, alias := range cmd.Aliases() {
			if commands.Registry[alias] != cmd {
				t.Errorf("Alias '%s' não aponta para o mesmo comando registrado", alias)
			}
		}
	}
}
