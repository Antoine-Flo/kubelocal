package command

import (
	"os/exec"
	"strings"

	"github.com/Antoine-Flo/kubelocal/internal/logger"
	"go.uber.org/zap"
)

func Run(log *logger.Logger, args ...string) error {
	_, err := RunWithOutput(log, args...)
	return err
}

func RunWithOutput(log *logger.Logger, args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	log.Debug("Executing command",
		zap.String("command", strings.Join(args, " ")))

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if len(outputStr) > 0 {
		log.Debug("Command output", zap.String("output", outputStr))
	}

	if err != nil {
		log.Debug("Command failed", zap.Error(err))
		return outputStr, err
	}

	return outputStr, nil
}
