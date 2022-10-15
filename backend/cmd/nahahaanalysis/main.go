package main

import (
	"time"

	"github.com/ohnishi/nahaha/backend/common/command"
	"github.com/spf13/cobra"
)

func transformRelateCommand() *cobra.Command {
	var (
		src   string
		dest  string
		dates []string
	)

	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Transform relate thread",
		RunE: command.WithLoggingE(func(cmd *cobra.Command, args []string) error {
			return command.EachDate(dates, func(date time.Time) error {
				return transformTrends(src, dest, date)
			})
		}),
	}
	cmd.PersistentFlags().StringVar(&src, "src", "~/Desktop", "dir to save spotify json")
	cmd.PersistentFlags().StringVar(&dest, "dest", "~/Desktop", "dir to save spotify json")
	command.SetDatesFlag(cmd.Flags(), &dates, "date for which the URL list file(s) is generated")
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{Use: "nahahaanalysis"}
	rootCmd.AddCommand(
		transformRelateCommand(),
	)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
