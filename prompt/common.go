package prompt

import (
	"log"

	"github.com/spf13/cobra"
)

func RootCmdNavigation(filename string, cmd *cobra.Command, args []string) {
	cmdName, err := SelectNextCmd(filename)
	if err != nil {
		log.Fatal(err)
	}

	nextCmd, _, err := cmd.Find([]string{cmdName})
	if err != nil {
		log.Fatal(err)
	}
	nextCmd.Run(nextCmd, args)
}
