package cmd

import (
	"context"
	"errors"
	"fmt"

	yogourtcompiler "github.com/goyourt/yogourt-compiler"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile tous les plugins Yogourt",
	Long:  "Rebuild tous les plugins Yogourt vers le dossier .yogourt.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, cfgPath, err := compilerConfig(nil)
		if err != nil {
			return err
		}

		fmt.Printf("🥛 Yogourt build\n")
		fmt.Printf("Config: %s\n", cfgPath)
		fmt.Printf("Routes: %s\n", cfg.RouteFolder)

		report, err := yogourtcompiler.Build(context.Background(), cfg)
		if err != nil {
			return err
		}

		for _, result := range report.Results {
			if result.Error != nil {
				fmt.Printf("❌ %s: %v\n", result.Source, result.Error)
				continue
			}
			fmt.Printf("✅ %s -> %s (%dms)\n", result.Source, result.Plugin, result.DurationMS)
		}
		fmt.Printf("Build finished success=%v duration=%s\n", report.Success, report.FinishedAt.Sub(report.StartedAt))

		if !report.Success {
			return errors.New("plugin build failed")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(BuildCmd)
}
