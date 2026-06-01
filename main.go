package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tq303/dedox/internal/document"
	"github.com/tq303/dedox/internal/parse"
)

var formats = []string{document.FileTypePdf, document.FileTypeDocx, document.FileTypeExcel, document.FileTypePowerPoint, document.FileTypeHtml, document.FileTypeRtf, document.FileTypeJpg}

var filterNames []string

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

		for _, name := range filterNames {
			filterFn, err := parse.Filters[name]
			if !err {
				log.Fatalf("unknown filter: %s (available: %s)", name, strings.Join(availableFilters(), ", "))
			}

			lines = filterFn(lines)
		}

		for _, line := range lines {
			fmt.Println(line)
		}
	},
}

func availableFilters() []string {
	names := make([]string, 0, len(parse.Filters))

	for key := range parse.Filters {
		names = append(names, key)
	}

	return names
}

func main() {
	rootCmd.Flags().StringArrayVarP(&filterNames, "filter", "f", nil, "apply a named filter (repeatable); available: pii, urls, ip, boilerplate, normalize, uniq")

	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
	}
}
