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
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(config settings.Config) error {
	router := gin.Default()
	netAddr := config.Network.ListenAddress + ":" + strconv.Itoa(config.Network.ListenPort)

	s := &http.Server{
		Addr:           netAddr,
		Handler:        router,
		ReadTimeout:    2 * time.Second,                                             // this works because constants have an adaptive type
		WriteTimeout:   time.Duration(config.Network.ResponseTimeout) * time.Second, // here we have to do a type conversion
		MaxHeaderBytes: 1 << 20,
	}

	router.POST("/g10k", authMiddleware(runG10K, config.API.AuthMode, config.API.Token))

	if err := serve(s, config); err != nil {
		return fmt.Errorf("Could not start router %s", err)
	}

	return nil
}

func serve(s *http.Server, c settings.Config) error {
	e := c.Network.EnableTLS
	crt := c.Network.TLSCertFile
	k := c.Network.TLSKeyFile

	if !e {
		fmt.Printf("Serving the API via HTTP on %s\n", s.Addr)
		return s.ListenAndServe()
	} else {
		fmt.Printf("Serveing the API via HTTPS on %s\n", s.Addr)
		return s.ListenAndServeTLS(crt, k)
	}
}

func authMiddleware(fn gin.HandlerFunc, m bool, t string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authMode == false bypasses token validation check
		if m == false {
			fn(c)
		}

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

		if tokenString != t {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fn(c)
	}
}
