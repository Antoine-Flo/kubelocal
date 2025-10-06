package command

import (
	"os/exec"
	"strings"

	"github.com/Antoine-Flo/kubelocal/internal/logger"
	"go.uber.org/zap"
)

// Run executes a command with variadic arguments and logs its output using the provided logger
func Run(log *logger.Logger, args ...string) error {

	cmd := exec.Command(args[0], args[1:]...)
	log.Debug("Executing command",
		zap.String("command", strings.Join(args, " ")))

	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		log.Debug("Command output", zap.String("output", string(output)))
	}

	if err != nil {
		log.Debug("Command failed", zap.Error(err))
		return err
	}

	return nil
}
