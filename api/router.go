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
	"arrakis/settings"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "net/http"
)

func SetupRouter(c *settings.Config) error {
	router := gin.Default()

	router.POST("/g10k", testWebhook())

	if err := startRouter(router, c.ListenPort); err != nil {
		return fmt.Errorf("Could not start router %s", err)
	}

	return nil
}

func startRouter(r *gin.Engine, port int) error {
	return r.Run(":" + strconv.Itoa(port))
}

func testWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("Hurray the webhook works!")
	}
}
