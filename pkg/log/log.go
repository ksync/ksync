// The log package provides an abstraction for processing log messages with
// a variable config. Logged messages are eventually passed to standard out.
package log

import (
  "os"

  "github.com/sirupsen/logrus"
)

// Logger creates a new logger with the specified options
func Logger() *logrus.Logger {
  // TODO: add a real config
  // TODO: make configurable
  logger := logrus.New()
  logger.Out = os.Stdout

  return logger
}
