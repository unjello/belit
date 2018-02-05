// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "belit",
	Short: "Exteremely simple header-only package manager for C/C++",
	Long: `Bêlit is a simple CLI tool for managing header-only libraries
in C/C++.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
