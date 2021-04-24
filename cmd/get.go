package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"time"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get progress",
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetCmd(cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	},
}

func runGetCmd(cmd *cobra.Command, args []string) (err error) {
	if len(config.Repositories) <= 0 {
		return errors.New("リポジトリを指定してください")
	}
	for _, repository := range config.Repositories {
		fmt.Println(repository)
	}
	err = getProgress(config.Repositories[0], "masatora", "2021-04-25 00:00:00")
	return
}

// 指定した日に指定したユーザが行ったコミットを表示する
func getProgress(repository string, username string, date string) (err error) {
	layout := "2006-01-02 00:00:00"
	start_date, err := time.Parse(layout, date)
	end_date := start_date.AddDate(0, 0, 1)
	cmd := exec.Command(
		"git", "-C", repository, "log",
		"--oneline",
		"--reverse",
		"--author=" + username,
		"--since=" + start_date.Format(layout),
		"--until=" + end_date.Format(layout),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func init() {
	rootCmd.AddCommand(getCmd)
}
