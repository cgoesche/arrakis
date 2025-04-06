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
	"arrakis/internal/status"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Args:  cobra.MatchAll(cobra.OnlyValidArgs),
	Short: "Generate an API token",
	Long: `When the --auth or -a flag is used at the command line
API calls will require an API Token either.`,
	Run: func(cmd *cobra.Command, args []string) {
		genToken(conf.API.TokenHashAlgorithm)
	},
}

func genToken(a string) {
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
			err = fmt.Errorf("unknown hash algorithm '%s'", a)
		}
	}

	if err != nil {
		err = status.New(status.ErrTokenGen, "token generation failed", err)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(`Using %s to generate an API Token
Please store this in a safe location and do not share it with anyone

Token: %x
`, a, token)
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
		return "", fmt.Errorf("unknown hash SHA3 algorithm size")
	}

	h.Write([]byte(s))
	t := string(h.Sum(nil))

	return t, nil
}
