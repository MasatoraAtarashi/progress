package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
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

	var output string
	err = addMetaDate(&output)

	for _, repository := range config.Repositories {
		repository_name := strings.Split(repository, "/")
		output += "##" + repository_name[len(repository_name)-1] + "\n\n"

		var commits string
		commits, err = getProgress(repository, "masatora", "2021-04-25 00:00:00")
		output += commits
		output += "\n\n"
	}

	execLess(output)
	return
}

// 出力にメタデータを追加
func addMetaDate(output *string) (err error) {
	*output += "#" + "2021-04-25\n\n"
	return nil
}

// 指定した日に指定したユーザが行ったコミットを表示する
func getProgress(repository string, username string, date string) (output string, err error) {
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
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	output = string(out)
	return
}

// lessコマンドで表示
func execLess(str string) (err error) {
	cmd := exec.Command("less", "-R")
	cmd.Stdin = strings.NewReader(str)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return
}

func init() {
	rootCmd.AddCommand(getCmd)
}
