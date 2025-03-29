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

type Config struct {
	ListenAddress string `mapstructure:"address"`
	ListenPort    int    `mapstructure:"port"`
	AuthMode      bool   `mapstructure:"auth"`
	Token         string `mapstructure:"token"`
	DebugMode     bool   `mapstructure:"debug"`
}

func SetDefault() Config {
	return Config{
		ListenAddress: "127.0.0.1",
		ListenPort:    8080,
		AuthMode:      false,
		Token:         "",
		DebugMode:     false,
	}
}
