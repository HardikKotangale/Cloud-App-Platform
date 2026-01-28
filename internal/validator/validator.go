package validator

import (
	"fmt"
	"strings"

	"github.com/HardikKotangale/Cloud-App-Platform/internal/spec"
)

func Validate(a *spec.AppSpec) []string {
	var issues []string

	trim := func(s string) string { return strings.TrimSpace(s) }

	if trim(a.Name) == "" {
		issues = append(issues, "name is required")
	}
	if trim(a.Namespace) == "" {
		issues = append(issues, "namespace is required")
	}
	if trim(a.Image) == "" {
		issues = append(issues, "image is required")
	}

	// Security + stability policy: forbid :latest
	if strings.HasSuffix(a.Image, ":latest") || (!strings.Contains(a.Image, ":") && a.Image != "") {
		issues = append(issues, fmt.Sprintf("image tag must be pinned (':latest' is not allowed): %q", a.Image))
	}

	// Operational policy: replicas bounds
	if a.Replicas < 1 || a.Replicas > 10 {
		issues = append(issues, fmt.Sprintf("replicas must be between 1 and 10 (got %d)", a.Replicas))
	}

	// Operational policy: port sanity
	if a.Port <= 0 || a.Port > 65535 {
		issues = append(issues, fmt.Sprintf("port must be 1-65535 (got %d)", a.Port))
	}

	// Require resource requests/limits (basic governance)
	if trim(a.Resources.CPU) == "" {
		issues = append(issues, "resources.cpu is required (e.g., \"250m\")")
	}
	if trim(a.Resources.Memory) == "" {
		issues = append(issues, "resources.memory is required (e.g., \"256Mi\")")
	}

	return issues
}
