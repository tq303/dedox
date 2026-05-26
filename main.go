package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tq303/ddx/internal/document"
	"github.com/tq303/ddx/internal/parse"
)

var formats = []string{document.FileTypePdf, document.FileTypeDocx, document.FileTypeExcel, document.FileTypePowerPoint}

var rootCmd = &cobra.Command{
	Use:   "ddx [file]",
	Short: "CLI file parser",
	Long:  "Parse and output text from supported file types: " + strings.Join(formats, ", "),
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

		parse.Filter(lines)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
