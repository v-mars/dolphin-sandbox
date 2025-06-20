package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var (
	Ver       = "v1.0.1"
	BuildDate = time.Now().Format("2006.01.02")
)

func Version() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of aigc",
		Long:  `This is dolphin aigc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version   : %s\n", Ver)
			fmt.Printf("Commit    : %s\n", "")
			fmt.Printf("BuildDate : %s\n", BuildDate)
		},
	}
	return versionCmd
}
