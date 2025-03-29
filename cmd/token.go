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
	"fmt"

	"crypto/sha256"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 3), cobra.OnlyValidArgs),
	Short: "Generate an API token",
	Long: `When the --auth, -a is used at the command line
API calls will require an API Token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := genToken(); err != nil {
			return fmt.Errorf("Error: %s", err)
		}
		return nil
	},
}

func genToken() error {
	s := "Random text to hash"
	h := sha256.New()

	h.Write([]byte(s))

	bs := h.Sum(nil)
	fmt.Printf(`Please store this in a safe location and do not share with anyone\n
	API Token: %x`, bs)

	return nil
}
