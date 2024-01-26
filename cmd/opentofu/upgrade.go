package main

import (
	"github.com/crosbygw/opentofu/pkg/opentofu"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *opentofu.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Install(cmd.Context())
		},
	}
	return cmd
}
