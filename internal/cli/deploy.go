package cli

import (
	"fmt"

	"github.com/HardikKotangale/Cloud-App-Platform/internal/kube"
	"github.com/HardikKotangale/Cloud-App-Platform/internal/render"
	"github.com/HardikKotangale/Cloud-App-Platform/internal/spec"
	"github.com/HardikKotangale/Cloud-App-Platform/internal/validator"
	"github.com/HardikKotangale/Cloud-App-Platform/internal/observability"
	"github.com/spf13/cobra"
)

func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <app.yaml>",
		Short: "Deploy an app to Kubernetes (lifecycle automation)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			app, err := spec.LoadFromFile(path)
			if err != nil {
				return err
			}

			issues := validator.Validate(app)
			if len(issues) > 0 {
				_ = observability.IncPolicyViolations(int64(len(issues)))
				_ = observability.IncDeployFailure()
				fmt.Println("❌ Validation failed:")
				for _, is := range issues {
					fmt.Printf("  - %s\n", is)
				}
				return fmt.Errorf("deploy blocked: spec violates governance rules")
			}

			// Ensure namespace exists (governance baseline)
			if err := kube.EnsureNamespace(app.Namespace); err != nil {
				return err
			}

			manifests, err := render.RenderManifests(app)
			if err != nil {
				return err
			}

			// Apply manifests
			if err := kube.ApplyYAML(manifests); err != nil {
				return err
			}
			_ = observability.IncDeploySuccess()
			fmt.Printf("✅ Deployed %q into namespace %q\n", app.Name, app.Namespace)
			fmt.Println("Tip: run `platformctl status -n", app.Namespace, app.Name, "`")
			return nil
		},
	}

	return cmd
}
