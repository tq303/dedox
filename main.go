package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tq303/ddx/internal/file"
)

var rootCmd = &cobra.Command{
	Use:   "ddx [file]",
	Short: "CLI file parser",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file path is required, usage: ddx [file]")
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		lines, err := file.Read(args[0])
		if err != nil {
			return
		}

		file.Filter(lines)
	},
}

func main() {
	rootCmd.Execute()
}
