package observability

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Counters struct {
	DeploySuccess    int64 `json:"deploy_success"`
	DeployFailure    int64 `json:"deploy_failure"`
	PolicyViolations int64 `json:"policy_violations"`
}

var mu sync.Mutex

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	dir := filepath.Join(home, ".platformctl")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return filepath.Join(dir, "metrics.json"), nil
}

func Load() (Counters, error) {
	mu.Lock()
	defer mu.Unlock()

	path, err := stateFilePath()
	if err != nil {
		return Counters{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Counters{}, nil
		}
		return Counters{}, fmt.Errorf("read metrics: %w", err)
	}

	var c Counters
	if err := json.Unmarshal(b, &c); err != nil {
		return Counters{}, fmt.Errorf("parse metrics json: %w", err)
	}
	return c, nil
}

func Save(c Counters) error {
	mu.Lock()
	defer mu.Unlock()

	path, err := stateFilePath()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metrics: %w", err)
	}

	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write metrics: %w", err)
	}
	return nil
}

func IncDeploySuccess() error {
	c, err := Load()
	if err != nil {
		return err
	}
	c.DeploySuccess++
	return Save(c)
}

func IncDeployFailure() error {
	c, err := Load()
	if err != nil {
		return err
	}
	c.DeployFailure++
	return Save(c)
}

func IncPolicyViolations(n int64) error {
	c, err := Load()
	if err != nil {
		return err
	}
	c.PolicyViolations += n
	return Save(c)
}
