/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package chat

import (
	"fmt"
	"os"
	"strings"

	"github.com/namnd/xai-cli/local"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "List chat history",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		history, err := local.ListRecentChatThread()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for i, h := range history {
			fmt.Println(fmt.Sprintf("%d)", i), strings.ReplaceAll(h.Prompt, "\n", ", "))
		}
	},
}

func init() {
	ChatCmd.AddCommand(historyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// historyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// historyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
