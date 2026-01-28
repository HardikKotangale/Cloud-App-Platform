package cli

import (
	"fmt"

	"github.com/HardikKotangale/cloud-app-platform/internal/kube"
	"github.com/HardikKotangale/cloud-app-platform/internal/render"
	"github.com/HardikKotangale/cloud-app-platform/internal/spec"
	"github.com/HardikKotangale/cloud-app-platform/internal/validator"
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

			fmt.Printf("✅ Deployed %q into namespace %q\n", app.Name, app.Namespace)
			fmt.Println("Tip: run `platformctl status -n", app.Namespace, app.Name, "`")
			return nil
		},
	}

	return cmd
}
