package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"dev-metrics/internal/commands"
	_ "dev-metrics/internal/commands" // O "_" dispara os inits e registra os comandos
)

func main() {
	flag.Usage = func() {
		fmt.Println("BMT - Build Metric Tool")
		fmt.Println("\nComandos disponíveis:")

		// Ordena os nomes para exibição bonita
		keys := make([]string, 0, len(commands.Registry))
		for k := range commands.Registry {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			cmd := commands.Registry[k]
			fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Description())
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
