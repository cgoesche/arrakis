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
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os/exec"
)

func runG10K(c *gin.Context) {
	cmd := exec.Command("/usr/bin/g10k", "-config", "/etc/puppetlabs/g10k/g10k.yaml")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed to attach to STDOUT",
			"error":  err})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed to attach to STDOUT",
			"error":  err})
		return
	}

	if err = cmd.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed to start program",
			"error":  err})
		return
	}

	fmt.Printf("Started g10k ...\n")
	stdOut, _ := io.ReadAll(stdout)
	stdErr, _ := io.ReadAll(stderr)

	if err = cmd.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "The program exited with a non-zero code",
			"error":  err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"output": string(stdOut) + string(stdErr),
		"error":  "",
	})
	return
}
