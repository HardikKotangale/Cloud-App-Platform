package cli

import (
	"fmt"

	"github.com/HardikKotangale/cloud-app-platform/internal/kube"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "status <app-name>",
		Short: "Show app status + diagnostics (pods, service, events)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			if namespace == "" {
				return fmt.Errorf("namespace is required: use -n <namespace>")
			}

			fmt.Printf("📦 App: %s (namespace: %s)\n\n", appName, namespace)

			fmt.Println("== Pods ==")
			if err := kube.PrintPods(namespace, appName); err != nil {
				return err
			}

			fmt.Println("\n== Service ==")
			if err := kube.PrintService(namespace, appName); err != nil {
				return err
			}

			fmt.Println("\n== Recent Events (Diagnostics) ==")
			if err := kube.PrintEvents(namespace); err != nil {
				return err
			}

			fmt.Println("\n✅ Done.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	return cmd
}
