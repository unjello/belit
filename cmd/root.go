// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var Verbose bool
var Noisy bool
var Debug bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output (info)")
	rootCmd.PersistentFlags().BoolVarP(&Noisy, "noisy", "n", false, "noisy output (trace)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "D", false, "debug output ( w/ output from commands)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("noisy", rootCmd.PersistentFlags().Lookup("noisy"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("verbose", false)
	viper.SetDefault("noisy", false)
	viper.SetDefault("debug", false)

	viper.BindEnv("cxx", "CXX")
	switch runtime.GOOS {
	case "linux":
		viper.SetDefault("cxx", "/usr/bin/g++")
		break
	case "darwin":
		viper.SetDefault("cxx", "/usr/bin/clang++")
		break
	case "windows":
		log.WithFields(log.Fields{
			"os": runtime.GOOS,
		}).Panic("System not supported")
	}
}

func initConfig() {
	if viper.GetBool("verbose") == true {
		logrus.SetLevel(logrus.InfoLevel)
	} else if viper.GetBool("noisy") == true || viper.GetBool("debug") == true {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	log.WithFields(log.Fields{
		"verbose": viper.GetBool("verbose"),
	}).Debug("Config value set")
	log.WithFields(log.Fields{
		"noisy": viper.GetBool("noisy"),
	}).Debug("Config value set")
	log.WithFields(log.Fields{
		"debug": viper.GetBool("debug"),
	}).Debug("Config value set")
}
