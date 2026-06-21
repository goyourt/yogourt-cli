package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yogourtcompiler "github.com/goyourt/yogourt-compiler"
	"github.com/spf13/cobra"
)

var startRuntimeCommand string

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Lance le runtime Yogourt",
	Long:  "Lance le runtime web Yogourt sans compiler les plugins et sans watcher.",
	RunE: func(cmd *cobra.Command, args []string) error {
		runtimeCommand := yogourtcompiler.GoRunCommand(".")
		if startRuntimeCommand != "" {
			runtimeCommand = strings.Fields(startRuntimeCommand)
		}
		if len(args) > 0 {
			runtimeCommand = args
		}
		if len(runtimeCommand) == 0 {
			return errors.New("commande runtime vide")
		}

		cfg, cfgPath, err := compilerConfig(runtimeCommand)
		if err != nil {
			return err
		}

		fmt.Printf("🥛 Yogourt start\n")
		fmt.Printf("Config: %s\n", cfgPath)
		fmt.Printf("Routes: %s\n", cfg.RouteFolder)
		fmt.Printf("Runtime: %s\n", strings.Join(runtimeCommand, " "))

		runtime := exec.Command(runtimeCommand[0], runtimeCommand[1:]...)
		runtime.Dir = "."
		runtime.Stdout = os.Stdout
		runtime.Stderr = os.Stderr
		runtime.Stdin = os.Stdin

		return runtime.Run()
	},
}

func init() {
	StartCmd.Flags().StringVar(&startRuntimeCommand, "runtime", "", "Commande runtime à lancer, par défaut: go run .")
	rootCmd.AddCommand(StartCmd)
}
