/*
Copyright © 2025 Christian Goeschel Ndjomouo <cgoesc2@wgu.edu>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"arrakis/app"
	"arrakis/settings"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile          string
	netAddr             string
	port                int
	responseTimeout     int
	enableTLS           bool
	tlsKeyFile          string
	tlsCertFile         string
	authMode            bool
	token               string
	tokenHashAlgo       string
	printHashAlgorithms bool
	debugMode           bool

	config settings.Config

	rootCmd = &cobra.Command{
		Use:   "arrakis",
		Short: "A lightweight git repository webhook API server",
		Long: `Although Arrakis is first and foremost the Desert Planet in the Dune universe, 
it is also a lightweight webhook API server that aims to handle webhooks triggered by 
arbitrary Git repository events and, depending on the payload, perform user-defined actions.`,
		Version: app.Version,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", settings.SetDefault().Logging.DebugMode, "enable verbose output for debugging")
	serverCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file to use")
	serverCmd.PersistentFlags().StringVarP(&netAddr, "address", "a", settings.SetDefault().Network.ListenAddress, "network address (e.g. 0.0.0.0)")
	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", settings.SetDefault().Network.ListenPort, "port to listen on")
	serverCmd.PersistentFlags().IntVar(&responseTimeout, "timeout", settings.SetDefault().Network.ResponseTimeout, "API request response timeout (sec)")
	serverCmd.PersistentFlags().BoolVarP(&enableTLS, "tls", "s", settings.SetDefault().Network.EnableTLS, "enable TLS")
	serverCmd.PersistentFlags().StringVar(&tlsCertFile, "cert", settings.SetDefault().Network.TLSCertFile, "TLS cert file path")
	serverCmd.PersistentFlags().StringVar(&tlsKeyFile, "key", settings.SetDefault().Network.TLSKeyFile, "TLS key file path")
	serverCmd.PersistentFlags().BoolVar(&authMode, "auth", settings.SetDefault().API.AuthMode, "enable verbose output for debugging")
	serverCmd.PersistentFlags().StringVarP(&token, "token", "t", settings.SetDefault().API.Token, "API token (implies --auth)")
	tokenCmd.PersistentFlags().StringVarP(&tokenHashAlgo, "algorithm", "a", settings.SetDefault().API.TokenHashAlgorithm, "specify the token hashing algorithm")

	viper.BindPFlag("logging.debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("network.address", serverCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("network.port", serverCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("network.responseTimeout", serverCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("network.enableTLS", serverCmd.PersistentFlags().Lookup("tls"))
	viper.BindPFlag("network.tlsCert", serverCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("network.tlsKey", serverCmd.PersistentFlags().Lookup("cert"))
	viper.BindPFlag("api.authMode", serverCmd.PersistentFlags().Lookup("auth"))
	viper.BindPFlag("api.token", serverCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("api.tokenHashAlgorithm", tokenCmd.PersistentFlags().Lookup("algorithm"))

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(versionCmd)

}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to find user's config directory!")
			os.Exit(1)
		}

		var configPath string
		configPath = filepath.Join(configDir, app.Name)

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix(app.Name)
	viper.AutomaticEnv()
	// Needed so that the viper engine can map the right suboptions
	// from the YAML configuration to the env
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("Unable to read config file: %s\nError: %s\n", viper.ConfigFileUsed(), err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to read values into configuration!")
		os.Exit(1)
	}

	if authMode && config.API.Token == "" {
		fmt.Printf("API token missing for auth mode. Please ")
		os.Exit(2)
	} else if !authMode && config.API.Token != "" {
		fmt.Printf("You have to enable authMode when using an API token")
		os.Exit(2)
	}
}
