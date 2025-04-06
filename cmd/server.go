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
	"arrakis/api"
	"arrakis/internal/status"
	"os"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve", "api", "listen"},
	Args:    cobra.MatchAll(cobra.OnlyValidArgs),
	Short:   "Start the API server",
	Long:    "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		serverStart()
	},
}

func serverStart() {
	conf.Logging.Logger.Debug("Starting the API server")

	if err := api.Start(conf); err != nil {
		err = status.New(status.ErrAPIStart, "unable to start API", err)
		conf.Logging.Logger.Error(err)
		os.Exit(1)
	}
}
