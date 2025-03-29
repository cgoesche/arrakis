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
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long: `There is not much more to say about this or 
maybe you are looking for the entire commit history ?`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := printVersion(app.Version); err != nil {
			return err
		}
		return nil
	},
}

func printVersion(v string) error {
	if _, err := fmt.Printf("arrakis version %s\n", v); err != nil {
		return err
	}
	return nil
}
