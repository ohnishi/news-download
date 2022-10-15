package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ohnishi/nahaha/backend/common/command"
	"github.com/spf13/cobra"
)

func main() {
	var (
		dates []string
		src   string
		dest  string
	)
	cmdFanza := &cobra.Command{
		Use:   "trends",
		Short: "Publish trends",
		Long:  "Publish trends from json",
		Args:  cobra.NoArgs,
		RunE: command.WithLoggingE(func(cmd *cobra.Command, args []string) error {
			return command.EachDate(dates, func(date time.Time) error {
				return publishTrends(src, dest, date)
			})
		}),
	}
	command.SetDatesFlag(cmdFanza.Flags(), &dates, "date for which the URL list file(s) is generated")
	_ = cmdFanza.MarkFlagRequired("date")
	cmdFanza.Flags().StringVar(&src, "src", "fanza/transform", "output path into which 5ch threads is written.")
	cmdFanza.Flags().StringVar(&dest, "dest", "./hugo/content/posts", "output path into which 5ch threads is written.")

	rootCmd := &cobra.Command{Use: "nahahapublish"}
	rootCmd.AddCommand(
		cmdFanza,
	)

	fmt.Printf("> %s\n", strings.Join(os.Args, " "))

	err := rootCmd.Execute()
	if err != nil {
		command.PrintErrorAndExit(err)
	}
}
