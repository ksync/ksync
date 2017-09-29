
package log

import (
  "os"

  "github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
  // TODO: add a real config
  // TODO: make configurable
  logger := logrus.New()
  logger.Out = os.Stdout

  return logger
}
