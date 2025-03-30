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
package settings

type Network struct {
	ListenAddress string `mapstructure:"address"`
	ListenPort    int    `mapstructure:"port"`
	EnableTLS     bool   `mapstructure:"enableTLS"`
}

type API struct {
	AuthMode           bool   `mapstructure:"authMode"`
	Token              string `mapstructure:"token"`
	TokenHashAlgorithm string `mapstructure:"tokenHashAlgorithm"`
}

type Logging struct {
	DebugMode bool `mapstructure:"debug"`
}

type Config struct {
	Network Network `mapstructure:"network"`
	API     API     `mapstructure:"api"`
	Logging Logging `mapstructure:"logging"`
}

func SetDefault() Config {
	return Config{
		Network: Network{
			ListenAddress: "127.0.0.1",
			ListenPort:    8080,
			EnableTLS:     false,
		},
		API: API{
			AuthMode:           false,
			Token:              "",
			TokenHashAlgorithm: "sha256",
		},
		Logging: Logging{
			DebugMode: false,
		},
	}
}
