package cli

import (
	"fmt"

	"github.com/HardikKotangale/cloud-app-platform/internal/spec"
	"github.com/HardikKotangale/cloud-app-platform/internal/validator"
	"github.com/spf13/cobra"
)

func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <app.yaml>",
		Short: "Validate an application spec (governance & security checks)",
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
				return fmt.Errorf("validation failed (%d issue(s))", len(issues))
			}

			fmt.Println("✅ Validation passed.")
			return nil
		},
	}

	return cmd
}
