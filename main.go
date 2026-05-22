package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tq303/ddx/internal/document"
)

var rootCmd = &cobra.Command{
	Use:   "ddx [file]",
	Short: "CLI file parser",
	Long:  "Parse and output text from supported file types: .txt, .pdf",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file path is required, usage: ddx [file]")
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		lines, err := document.Read(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		document.Filter(lines)
	},
}

func main() {
	rootCmd.Execute()
}
