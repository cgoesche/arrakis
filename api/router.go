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
	"strings"

	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(config *settings.Config) error {
	router := gin.Default()
	netAddr := config.ListenAddress + ":" + strconv.Itoa(config.ListenPort)

	router.POST("/g10k", authWebhook(&config.Token))

	if err := startRouter(router, netAddr); err != nil {
		return fmt.Errorf("Could not start router %s", err)
	}

	return nil
}

func startRouter(r *gin.Engine, a string) error {
	return r.Run(a)
}

func authWebhook(t *string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authToken[1]
		if tokenString != *(t) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
