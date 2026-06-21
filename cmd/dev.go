package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goyourt/yogourt-cli/config"
	yogourtcompiler "github.com/goyourt/yogourt-compiler"
	"github.com/spf13/cobra"
)

var devRuntimeCommand string
var devNoInitialBuild bool

var DevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Lance Yogourt en mode développement",
	Long:  "Compile les plugins Yogourt, lance le runtime web, watch les changements, puis rebuild/restart en développement.",
	RunE: func(cmd *cobra.Command, args []string) error {
		runtimeCommand := yogourtcompiler.GoRunCommand(".")
		if devRuntimeCommand != "" {
			runtimeCommand = strings.Fields(devRuntimeCommand)
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

		fmt.Printf("🥛 Yogourt dev\n")
		fmt.Printf("Config: %s\n", cfgPath)
		fmt.Printf("Routes: %s\n", cfg.RouteFolder)
		fmt.Printf("Runtime: %s\n", strings.Join(runtimeCommand, " "))
		if cfg.NoInitialBuild {
			fmt.Println("Initial build: skipped")
		}

		return yogourtcompiler.Dev(context.Background(), cfg)
	},
}

func compilerConfig(runtimeCommand []string) (yogourtcompiler.Config, string, error) {
	cfg, cfgPath, err := loadProjectConfig()
	if err != nil {
		return yogourtcompiler.Config{}, "", err
	}

	routeFolder := cfg.Paths.APIFolder
	if routeFolder == "" {
		routeFolder = cfg.Paths.RouteFolder
	}
	if routeFolder == "" {
		routeFolder = "api"
	}
	routeFolder = strings.Trim(routeFolder, "./")

	return yogourtcompiler.Config{
		Root:           ".",
		RouteFolder:    routeFolder,
		MiddlewareFile: "middleware/middleware.go",
		RuntimeCommand: runtimeCommand,
		NoInitialBuild: devNoInitialBuild,
		WatchPaths: []string{
			routeFolder,
			"middleware",
			"models",
			"controllers",
			"bridge",
			"services",
			"cmd",
			"go.mod",
			"go.sum",
		},
	}, cfgPath, nil
}

func loadProjectConfig() (*config.Config, string, error) {
	candidates := []string{
		filepath.Join("configs", "yogourt.yaml"),
		"config.yaml",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			cfg, err := config.LoadConfig(candidate)
			return cfg, candidate, err
		}
	}

	return nil, "", fmt.Errorf("❌ Fichier config Yogourt introuvable: attendu %s ou %s", candidates[0], candidates[1])
}

func init() {
	DevCmd.Flags().StringVar(&devRuntimeCommand, "runtime", "", "Commande runtime à lancer, par défaut: go run .")
	DevCmd.Flags().BoolVar(&devNoInitialBuild, "no-initial-build", false, "Lance le runtime sans rebuild initial des plugins")
	rootCmd.AddCommand(DevCmd)
}
