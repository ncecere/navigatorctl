package cmd

import (
	"fmt"
	"os"

	"github.com/ncecere/navigatorctl/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "navigatorctl",
		Short: "A CLI tool for managing Navigator resources",
		Long: `navigatorctl is a command line interface for managing Navigator resources.
It provides functionality for managing teams, members, and API keys.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.navigatorctl.yaml)")

	// API configuration flags
	rootCmd.PersistentFlags().String("api-url", "", "API URL")
	rootCmd.PersistentFlags().String("api-key", "", "API Key")

	// Bind flags to viper
	viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("api.key", rootCmd.PersistentFlags().Lookup("api-key"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".navigatorctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".navigatorctl")

		// Also look for config in the current directory
		viper.AddConfigPath(".")
	}

	// Read in environment variables that match
	viper.SetEnvPrefix("NAVIGATOR")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}

	// Validate required configuration
	if viper.GetString("api.url") == "" {
		fmt.Fprintln(os.Stderr, "Error: API URL is required. Set it in config file or use --api-url flag")
		os.Exit(1)
	}

	if viper.GetString("api.key") == "" {
		fmt.Fprintln(os.Stderr, "Error: API key is required. Set it in config file or use --api-key flag")
		os.Exit(1)
	}
}

// getAPIClient creates a new API client using the current configuration
func getAPIClient() *api.Client {
	return api.NewClient(
		viper.GetString("api.url"),
		viper.GetString("api.key"),
	)
}
