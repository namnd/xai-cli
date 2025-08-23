/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package chat

import (
	"fmt"
	"os"
	"strconv"

	"github.com/namnd/xai-cli/local"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume [number]",
	Short: "Resume a chat",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = append(args, "0")
		}

		t, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		chat, err := local.GetChatMinusT(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		chat.Display()
	},
}

func init() {
	ChatCmd.AddCommand(resumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
