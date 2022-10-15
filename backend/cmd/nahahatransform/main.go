package main

import (
	"time"

	"github.com/ohnishi/nahaha/backend/common/command"
	"github.com/spf13/cobra"
)

func transformRSSCommand() *cobra.Command {
	var (
		src   string
		dest  string
		dates []string
	)

	cmd := &cobra.Command{
		Use:   "rss",
		Short: "Transform rss thread",
		RunE: command.WithLoggingE(func(cmd *cobra.Command, args []string) error {
			return command.EachDate(dates, func(date time.Time) error {
				return transformRSS(src, dest, date)
			})
		}),
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	return command.EachDate(dates, func(date time.Time) error {
		// 		return transformNews5ch(src, dest, date)
		// 	})
		// },
	}
	cmd.PersistentFlags().StringVar(&src, "src", "~/Desktop", "dir to save spotify json")
	cmd.PersistentFlags().StringVar(&dest, "dest", "~/Desktop", "dir to save spotify json")
	command.SetDatesFlag(cmd.Flags(), &dates, "date for which the URL list file(s) is generated")
	_ = cmd.MarkFlagRequired("date")

	return cmd
}

func transform5chCommand() *cobra.Command {
	var (
		src   string
		dest  string
		dates []string
	)

	cmd := &cobra.Command{
		Use:   "5ch",
		Short: "Transform 5ch thread",
		RunE: command.WithLoggingE(func(cmd *cobra.Command, args []string) error {
			return command.EachDate(dates, func(date time.Time) error {
				return transform5ch(src, dest, date)
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
	rootCmd := &cobra.Command{Use: "nahahatransform"}
	rootCmd.AddCommand(
		transformRSSCommand(),
		transform5chCommand(),
	)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
