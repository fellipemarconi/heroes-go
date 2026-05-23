package containers

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Summary struct {
	Name   string
	Status string
}

const (
	listTimeout    = 5 * time.Second
	restartTimeout = 30 * time.Second
)

func List(ctx context.Context) ([]Summary, error) {
	output, err := runDocker(ctx, listTimeout, "ps", "-a", "--format", "{{.Names}}\t{{.Status}}")
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []Summary{}, nil
	}

	lines := strings.Split(trimmed, "\n")
	containers := make([]Summary, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected docker output: %q", line)
		}

		containers = append(containers, Summary{
			Name:   parts[0],
			Status: parts[1],
		})
	}

	return containers, nil
}

func RestartAll(ctx context.Context) error {
	runningIDs, err := listIDs(ctx, true)
	if err != nil {
		return err
	}

	if len(runningIDs) > 0 {
		if _, err := runDocker(ctx, restartTimeout, append([]string{"stop"}, runningIDs...)...); err != nil {
			return err
		}
	}

	allIDs, err := listIDs(ctx, false)
	if err != nil {
		return err
	}

	if len(allIDs) == 0 {
		return nil
	}

	if _, err := runDocker(ctx, restartTimeout, append([]string{"start"}, allIDs...)...); err != nil {
		return err
	}

	return nil
}

func listIDs(ctx context.Context, runningOnly bool) ([]string, error) {
	args := []string{"ps", "-aq"}
	if runningOnly {
		args = []string{"ps", "-q"}
	}

	output, err := runDocker(ctx, listTimeout, args...)
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(output)), nil
}

func runDocker(ctx context.Context, timeout time.Duration, args ...string) ([]byte, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker %s failed: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}

	return output, nil
}
