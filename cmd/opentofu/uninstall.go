package main

import (
	"github.com/crosbygw/opentofu/pkg/opentofu"
	"github.com/spf13/cobra"
)

func buildUninstallCommand(m *opentofu.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Execute the uninstall functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Install(cmd.Context())
		},
	}
	return cmd
}
