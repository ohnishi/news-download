package command

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
)

// NewFlagErrorf はflagErrorを返す。
// cobra.Commandのフラグチェックに引っかからないフラグエラーを返したい場合に利用する。
func NewFlagErrorf(msg string, args ...interface{}) error {
	return flagError{Message: msg, Args: args}
}

type flagError struct {
	Message string
	Args    []interface{}
}

func (e flagError) Error() string {
	return fmt.Sprintf(e.Message, e.Args...)
}

// IsFlagError は引数errがflagErrorならtrueを返す
func IsFlagError(err error) bool {
	_, ok := err.(flagError)
	return ok
}

// WithLogging はコマンド実行時にエラーが起きたらログを出力する
func WithLogging(fn func(*cobra.Command, []string) error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		withLogging(fn, cmd, args)
	}
}

// WithLoggingE はコマンド実行時にエラーが起きたらログを出力してエラーを返す
func WithLoggingE(fn func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return withLogging(fn, cmd, args)
	}
}

func withLogging(fn func(cmd *cobra.Command, args []string) error, cmd *cobra.Command, args []string) error {
	err := fn(cmd, args)
	if err == nil {
		return nil
	}

	if err == flag.ErrHelp {
		return cmd.Help()
	}

	cmd.Printf("Error: %+v\n", err)
	if IsFlagError(err) {
		cmd.Usage()
	}
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	return err
}
