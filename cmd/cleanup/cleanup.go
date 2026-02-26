package cleanup

import (
	"fmt"
	"strings"

	"github.com/openshift/backplane-tools/pkg/tools"
	"github.com/openshift/backplane-tools/pkg/utils"
	"github.com/spf13/cobra"
)

// Cmd returns the Command used to invoke the cleanup logic
func Cmd() *cobra.Command {
	toolNames := tools.Names()
	cleanupCmd := &cobra.Command{
		Use:       fmt.Sprintf("cleanup [all|%s]", strings.Join(toolNames, "|")),
		Aliases:   []string{"clean-up"},
		Args:      cobra.OnlyValidArgs,
		ValidArgs: append(toolNames, "all"),
		Short:     "Cleans up older versions of existing tool",
		Long:      "Cleans up older versions and keeps only the latest installed version of one or more tools from the provided list. It's valid to specify multiple tools: in this case, all tools provided will be cleaned up. If no specific tools are provided, all are cleaned up by default.",
		RunE: func(_ *cobra.Command, args []string) error {
			return Cleanup(args)
		},
	}
	return cleanupCmd
}

// Cleanup cleans up older versions of the tool(s) specified by the provided positional args
func Cleanup(args []string) error {
	var listTools []tools.Tool
	if len(args) == 0 || utils.Contains(args, "all") {
		// If user explicitly passes 'all' or doesn't specify which tools to clean up,
		// everything that's been installed locally will be cleaned up
		var err error
		listTools, err = tools.ListInstalled()
		if err != nil {
			return err
		}
	} else {
		// otherwise build the list verifying tool exist
		toolMap := tools.GetMap()

		listTools = []tools.Tool{}
		for _, toolName := range args {
			t, found := toolMap[toolName]
			if !found {
				return fmt.Errorf("failed to locate '%s' in list of supported tools", toolName)
			}
			listTools = append(listTools, t)
		}
	}

	fmt.Println("Cleaning up the following tools: ")

	err := tools.Cleanup(listTools)
	if err != nil {
		return fmt.Errorf("failed to clean up tools: %w", err)
	}
	return nil
}
