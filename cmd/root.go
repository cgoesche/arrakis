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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	netAddr    string
	port       int
	authMode   bool
	token      string
	debugMode  bool

	config settings.Config

	rootCmd = &cobra.Command{
		Use:   "arrakis",
		Short: "A lightweight git repository webhook API server",
		Long: `Although Arrakis is first and foremost the Desert Planet in the Dune universe, 
it is also a lightweight webhook API server that aims to handle webhooks triggered by 
arbitrary Git repository events and, depending on the payload, perform user-defined actions.`,
		Version: app.Version,
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

	serverCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file to use")
	serverCmd.PersistentFlags().StringVarP(&netAddr, "address", "a", settings.SetDefault().ListenAddress, "network address (e.g. 0.0.0.0)")
	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", settings.SetDefault().ListenPort, "port to listen on")
	serverCmd.PersistentFlags().BoolVar(&authMode, "auth", settings.SetDefault().AuthMode, "enable verbose output for debugging")
	serverCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "API token (implies --auth)")
	serverCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", settings.SetDefault().DebugMode, "enable verbose output for debugging")
	// serverCmd.MarkFlagsRequiredTogether("auth", "token")

	viper.BindPFlag("address", serverCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("debug", serverCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("auth", serverCmd.PersistentFlags().Lookup("auth"))
	viper.BindPFlag("token", serverCmd.PersistentFlags().Lookup("token"))

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
		configPath = filepath.Join(configDir, "arrakis")

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix(app.Name)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to read values into configuration!")
		os.Exit(1)
	}

	if authMode && config.Token == "" {
		fmt.Printf("Token missing for auth mode")
		os.Exit(2)
	}
}
