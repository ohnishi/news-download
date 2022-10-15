package main

import (
	"github.com/spf13/cobra"
)

func new5chFetchCommand() *cobra.Command {
	var (
		dest string
	)

	cmd := &cobra.Command{
		Use:   "5ch",
		Short: "Fetch 5ch thread",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fetchNews5ch(dest, 3)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&dest, "dest", "~/Desktop", "dir to save spotify json")

	return cmd
}

func newYahooFetchCommand() *cobra.Command {
	var (
		dest string
	)

	cmd := &cobra.Command{
		Use:   "yahoo",
		Short: "Fetch yahoo thread",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fetchYahooRSS(dest, 3)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&dest, "dest", "~/Desktop", "dir to save spotify json")

	return cmd
}

func newRSSFetchCommand() *cobra.Command {
	var (
		src  string
		dest string
	)

	cmd := &cobra.Command{
		Use:   "rss",
		Short: "Fetch rss thread",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := fetchNewsRSS(src, dest, 3)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&src, "src", "~/Desktop", "dir to save spotify json")
	cmd.PersistentFlags().StringVar(&dest, "dest", "~/Desktop", "dir to save spotify json")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{Use: "nahahafetch"}
	rootCmd.AddCommand(
		new5chFetchCommand(),
		newYahooFetchCommand(),
		newRSSFetchCommand(),
	)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
