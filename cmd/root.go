// andrzej lichnerowicz, unlicensed (~public domain)

package cmd

import (
	"runtime"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unjello/belit/config"
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
		config.GetConfig().Log.Panic(err)
	}
}

// Verbose is 1st least noisy debugging level for Belit
var Verbose bool

// Noisy is 2nd, most noisy debugging level for Belit
var Noisy bool

// Debug is 3rd debugging level, that combines Noisy of Belit, and enables additional
// logging in spawned programs.
var Debug bool

func init() {
	cobra.OnInitialize(initConfig)

	addBoolFlag(&Verbose, false, "verbose", "v", "verbose output (info)")
	addBoolFlag(&Noisy, false, "noisy", "n", "noisy output (trace)")
	addBoolFlag(&Debug, false, "debug", "D", "debug output ( w/ output from commands)")

	viper.BindEnv("cxx", "CXX")
	viper.BindEnv("cc", "CC")
	viper.BindEnv("cxxopts", "CXXOPTS")
	viper.BindEnv("ccopts", "CCOPTS")

	switch runtime.GOOS {
	case "linux":
		viper.SetDefault("cxx", "/usr/bin/g++")
		viper.SetDefault("cc", "/usr/bin/gcc")
		break
	case "darwin":
		viper.SetDefault("cxx", "/usr/bin/clang++")
		viper.SetDefault("cc", "/usr/bin/clang")
		break
	case "windows":
		log.WithFields(log.Fields{
			"os": runtime.GOOS,
		}).Panic("System not supported")
	}
}

func initConfig() {
	log := config.GetConfig().Log
	if viper.GetBool("verbose") == true {
		log.SetLevel(logrus.InfoLevel)
		log.Info("Enabled verbose logging level")
	} else if viper.GetBool("noisy") == true {
		log.SetLevel(logrus.DebugLevel)
		log.Info("Enabled noisy logging level")
	} else if viper.GetBool("debug") == true {
		log.SetLevel(logrus.DebugLevel)
		log.Info("Enabled debug logging level")
	} else {
		log.SetLevel(logrus.ErrorLevel)
	}
}

func addBoolFlag(p *bool, defaultValue bool, name string, shortHand string, description string) {
	rootCmd.PersistentFlags().BoolVarP(p, name, shortHand, defaultValue, description)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
	viper.SetDefault(name, defaultValue)
}
