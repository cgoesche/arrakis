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
	"arrakis/api/handler"
	"arrakis/api/middleware"
	"arrakis/app"
	"arrakis/internal/config"
	"arrakis/internal/status"

	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(conf config.Config) error {
	r := gin.Default()
	a := conf.Network.ListenAddress + ":" + strconv.Itoa(conf.Network.ListenPort)
	t := conf.API.Token
	m := conf.API.AuthMode

	s := &http.Server{
		Addr:           a,
		Handler:        r,
		ReadTimeout:    2 * time.Second,                                           // this works because constants have an adaptive type
		WriteTimeout:   time.Duration(conf.Network.ResponseTimeout) * time.Second, // here we have to do a type conversion
		MaxHeaderBytes: 1 << 20,
	}

	r.GET("/v1/info", getAPIInfo())
	r.GET("/v1/health", healthCheck())
	r.POST("/v1/webhook/g10k", middleware.AuthCall(m, t), handler.RunG10K())

	if conf.API.AuthMode {
		conf.Logging.Logger.Debug("Running in auth mode")
	}

	if err := serve(s, conf); err != nil {
		return err
	}

	return nil
}

func serve(s *http.Server, c config.Config) error {
	tls := c.Network.EnableTLS
	crt := c.Network.TLSCertFile
	k := c.Network.TLSKeyFile

	if !tls {
		c.Logging.Logger.Info(fmt.Sprintf("Serving the API via HTTP on %s\n", s.Addr))
		return s.ListenAndServe()
	} else {
		c.Logging.Logger.Info(fmt.Sprintf("Serving the API via HTTPS on %s\n", s.Addr))
		err := s.ListenAndServeTLS(crt, k)
		if err != nil {
			err = status.New(status.ErrTLSFailed, "failed to use TLS", err)
		}
		return err
	}
}

func getAPIInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    app.Name,
			"version": app.Version,
		})
	}
}

func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "running"})
	}
}
