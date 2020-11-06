package cmd

import (
	"github.com/aaomidi/no-biden-or-trump/telegram"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	RunE: func(cmd *cobra.Command, args []string) error {

		tg := telegram.New(viper.GetString("token"))

		if err := tg.Create(); err != nil {
			return errors.Wrap(err, "error creating telegram bot")
		}

		go tg.Start()

		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

		toTerminate := <-terminate
		for {
			if toTerminate != nil {
				tg.Stop()
				break
			}
		}

		return nil
	},
}
