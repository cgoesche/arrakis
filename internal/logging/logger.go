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
package logging

import (
	"fmt"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func New() *Logger {
	logger := slog.New(slog.Default().Handler())

	return &Logger{
		logger: logger,
	}
}

func SetLogLevel(lvl string) {
	switch lvl {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}

func (l *Logger) Debug(msg any) {
	m := convMsg(msg)
	l.logger.Debug(m)
}

func (l *Logger) Info(msg any) {
	m := convMsg(msg)
	l.logger.Info(m)
}

func (l *Logger) Warn(msg any) {
	m := convMsg(msg)
	l.logger.Warn(m)
}

func (l *Logger) Error(msg any) {
	m := convMsg(msg)
	l.logger.Error(m)
}

func convMsg(msg any) string {
	switch m := msg.(type) {
	case error:
		return m.Error()
	default:
		return fmt.Sprintf("%s", m)
	}
}
