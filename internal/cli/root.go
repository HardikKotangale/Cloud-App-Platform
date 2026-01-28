package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "platformctl",
		Short: "Cloud Application Lifecycle Governance Platform (Kubernetes)",
		Long:  "platformctl onboards, validates, deploys, and diagnoses cloud-native applications on Kubernetes.",
	}

	root.AddCommand(NewValidateCmd())
	root.AddCommand(NewDeployCmd())
	root.AddCommand(NewStatusCmd())
	root.AddCommand(NewMetricsCmd())
	return root
}
