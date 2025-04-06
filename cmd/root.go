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
	"arrakis/internal/config"
	"arrakis/internal/logging"
	"arrakis/internal/status"

	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile      string
	netAddr         string
	port            int
	responseTimeout int
	enableTLS       bool
	tlsKeyFile      string
	tlsCertFile     string
	authMode        bool
	token           string
	tokenHashAlgo   string
	debugMode       bool

	conf config.Config

	rootCmd = &cobra.Command{
		Use:   "arrakis",
		Short: "A lightweight Puppet g10k API server",
		Long: `Although Arrakis is first and foremost home to the Fremen, it is also a lightweight API server
that aims to handle Puppet control repo webhooks for g10k. It supports token-based authentication and HTTPS.`,
		Version: app.Version,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	serverCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", config.SetDefault().Logging.DebugMode, "enable verbose output for debugging")
	serverCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file to use")
	serverCmd.PersistentFlags().StringVarP(&netAddr, "address", "a", config.SetDefault().Network.ListenAddress, "network address (e.g. 0.0.0.0)")
	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", config.SetDefault().Network.ListenPort, "port to listen on")
	serverCmd.PersistentFlags().IntVar(&responseTimeout, "timeout", config.SetDefault().Network.ResponseTimeout, "API request response timeout (sec)")
	serverCmd.PersistentFlags().BoolVarP(&enableTLS, "tls", "s", config.SetDefault().Network.EnableTLS, "enable TLS")
	serverCmd.PersistentFlags().StringVar(&tlsCertFile, "cert", config.SetDefault().Network.TLSCertFile, "TLS cert file path")
	serverCmd.PersistentFlags().StringVar(&tlsKeyFile, "key", config.SetDefault().Network.TLSKeyFile, "TLS key file path")
	serverCmd.PersistentFlags().BoolVar(&authMode, "auth", config.SetDefault().API.AuthMode, "require an API token for protected calls")
	serverCmd.PersistentFlags().StringVarP(&token, "token", "t", config.SetDefault().API.Token, "API token (implies --auth)")
	tokenCmd.PersistentFlags().StringVarP(&tokenHashAlgo, "algorithm", "a", config.SetDefault().API.TokenHashAlgorithm, "specify the token hashing algorithm")

	viper.BindPFlag("logging.debug", serverCmd.PersistentFlags().Lookup("debug"))
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
	var logger = logging.New()

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			err = status.New(status.ErrConfigNotFound, "unable to find user's config directory", err)
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		var configPath = filepath.Join(configDir, app.Name)

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.SetEnvPrefix(app.Name)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		err = status.New(status.ErrConfigInvalid, "unable to read config file", err)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		err = status.New(status.ErrLoadConfigValues, "unable to load values into configuration", err)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}

	if conf.Logging.DebugMode {
		logging.SetLogLevel("debug")
	}
	conf.Logging.Logger = logger

	if authMode && conf.API.Token == "" {
		err := status.New(status.ErrAPITokenMissing, "missing API token for auth mode", nil)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	} else if !authMode && conf.API.Token != "" {
		err := status.New(status.ErrAPINotAuthMode, "not running in auth mode", nil)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
}
