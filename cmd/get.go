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

const layout = "2006-01-02 00:00:00"

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

	date, err := getDate(cmd)
	if err != nil {
		return
	}

	var output string
	err = addMetaDate(&output, date)

	for _, repository := range config.Repositories {
		repository_name := strings.Split(repository, "/")
		output += "## " + repository_name[len(repository_name)-1] + "\n"

		var commits string
		commits, err = getProgress(repository, "masatora", date)
		output += commits
		output += "\n"
	}

	fmt.Println(output)
	return
}

// 出力にメタデータを追加
func addMetaDate(output *string, datetime string) (err error) {
	date := strings.Split(datetime, " ")[0]
	*output += "# " + date + "\n\n"
	return nil
}

func getDate(cmd *cobra.Command) (date string, err error) {
	date, err = cmd.PersistentFlags().GetString("date")
	if err != nil {
		return
	}
	if date == "" {
		//out, err := exec.Command("git", "config", "user.name").Output()
		//if err != nil {
		//	return :err
		//}
		//date = string(out)
		date = time.Now().Format(layout)
	} else {
		date += " 00:00:00"
	}
	return
}

// 指定した日に指定したユーザが行ったコミットを表示する
func getProgress(repository string, username string, date string) (output string, err error) {
	start_date, err := time.Parse(layout, date)
	end_date := start_date.AddDate(0, 0, 1)
	cmd := exec.Command(
		"git", "-C", repository, "log",
		"--oneline",
		"--reverse",
		"--author=" + username,
		"--since=" + start_date.Format(layout),
		"--until=" + end_date.Format(layout),
		"--format= - %C(auto)%h%Creset %s",
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	output = string(out)
	return
}

func init() {
	getCmd.PersistentFlags().StringP("date", "d", "", "Specify date")

	rootCmd.AddCommand(getCmd)
}
