package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"dev-metrics/internal/commands"
	_ "dev-metrics/internal/commands" // O "_" dispara os inits e registra os comandos
)

func main() {
	flag.Usage = func() {
		fmt.Println("BMT - Build Metric Tool")
		fmt.Println("\nComandos disponíveis:")

		sortedCommands := make([]commands.Command, len(commands.CommandsList))
		copy(sortedCommands, commands.CommandsList)
		sort.Slice(sortedCommands, func(i, j int) bool {
			return sortedCommands[i].Name() < sortedCommands[j].Name()
		})

		for _, cmd := range sortedCommands {
			s := cmd.Name()
			if len(cmd.Aliases()) > 0 {
				s += " (" + strings.Join(cmd.Aliases(), ", ") + ")"
			}
			fmt.Printf("  %s:\n\t%s\n", s, cmd.Description())
		}
	}

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Busca o comando no registro
	cmdName := args[0]
	if cmd, ok := commands.Registry[cmdName]; ok {
		err := cmd.Run(args[1:]) // Passa apenas os argumentos restantes
		if err != nil {
			fmt.Printf("Erro ao executar comando '%s': %v\n", cmdName, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Erro: comando '%s' não encontrado.\n", cmdName)
		flag.Usage()
		os.Exit(1)
	}
}
