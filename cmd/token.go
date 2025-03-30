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
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"fmt"
	"github.com/spf13/cobra"
	"regexp"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 3), cobra.OnlyValidArgs),
	Short: "Generate an API token",
	Long: `When the --auth or -a flag is used at the command line
API calls will require an API Token either.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := genToken(config.API.TokenHashAlgorithm); err != nil {
			return fmt.Errorf("Error: %s", err)
		}
		return nil
	},
}

func genToken(a string) error {
	var (
		token string
		err   error
	)
	s := rand.Text()
	sha3Pat := regexp.MustCompile(`^sha3-.*$`)

	switch a {
	case "sha256":
		token, err = genSHA256(s)
	case "sha512":
		token, err = genSHA512(s)
	default:
		if sha3Pat.MatchString(a) {
			token, err = genSHA3(s, a)
		} else {
			return fmt.Errorf("Unknown hash algorithm %s\n", a)
		}
	}

	fmt.Printf(`Using %s to generate an API Token
Please store this in a safe location and do not share it with anyone

API Token:
%x
`, a, token)

	return err
}

func genSHA256(s string) (string, error) {
	h := sha256.New()
	h.Write([]byte(s))
	t := string(h.Sum(nil))

	return t, nil
}

func genSHA512(s string) (string, error) {
	h := sha512.New()
	h.Write([]byte(s))
	t := string(h.Sum(nil))

	return t, nil
}

func genSHA3(s string, a string) (string, error) {
	var h *sha3.SHA3

	switch a {
	case "sha3-224":
		h = sha3.New224()
	case "sha3-256":
		h = sha3.New256()
	case "sha3-384":
		h = sha3.New384()
	case "sha3-512":
		h = sha3.New512()
	default:
		return "", fmt.Errorf("Unknown hash SHA3 algorithm size!")
	}

	h.Write([]byte(s))
	t := string(h.Sum(nil))

	return t, nil
}
