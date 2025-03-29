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
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve", "api", "listen"},
	Args:    cobra.MatchAll(cobra.RangeArgs(0, 3), cobra.OnlyValidArgs),
	Short:   "Start the API server",
	Long:    "Start the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serverStart(); err != nil {
			return fmt.Errorf("Error: %s", err)
		}
		return nil
	},
}

func serverStart() error {
	if err := api.SetupRouter(&config); err != nil {
		return err
	}

	fmt.Printf("Listening on Port: %d\n", config.ListenPort)

	return nil
}
