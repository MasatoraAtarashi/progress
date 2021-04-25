package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const layout = "2006-01-02 00:00:00"

type Commits struct {
	Content string
	Count   int
}

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

	username, err := getUserName(cmd)
	if err != nil {
		return
	}

	var output string
	err = addMetaDate(&output, date)

	for _, repository := range config.Repositories {
		var commits Commits
		commits, err = getProgress(repository, username, date)
		if commits.Count > 0 {
			repository_name := strings.Split(repository, "/")
			output += "## " + repository_name[len(repository_name)-1] + "(" + strconv.Itoa(commits.Count) + " commits)" + "\n"
			output += commits.Content
			output += "\n"
		}
	}

	fmt.Printf(output)
	return
}

// 出力にメタデータを追加
func addMetaDate(output *string, datetime string) (err error) {
	date := strings.Split(datetime, " ")[0]
	*output += "# " + date + "\n\n"
	return nil
}

// 日付を取得
func getDate(cmd *cobra.Command) (date string, err error) {
	date, err = cmd.PersistentFlags().GetString("date")
	if err != nil {
		return
	}
	if date == "" {
		date = time.Now().Format(layout)
	} else {
		date += " 00:00:00"
	}
	return
}

// usernameを取得
func getUserName(cmd *cobra.Command) (username string, err error) {
	username, err = cmd.PersistentFlags().GetString("user")
	if err != nil {
		return
	}

	if username == "" {
		out, err := exec.Command("git", "config", "user.name").Output()
		if err != nil {
			return "", err
		}
		username = string(out)
	}
	return
}

// 指定した日に指定したユーザが行ったコミットを表示する
func getProgress(repository string, username string, date string) (commit Commits, err error) {
	start_date, err := time.Parse(layout, date)
	end_date := start_date.AddDate(0, 0, 1)
	cmd := exec.Command(
		"git", "-C", repository, "log",
		"--oneline",
		"--author="+username,
		"--since="+start_date.Format(layout),
		"--until="+end_date.Format(layout),
		"--format= - %C(auto)%h%Creset %s",
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	commit.Content = string(out)
	commit.Count = len(strings.Split(commit.Content, "\n")) - 1
	return
}

func init() {
	getCmd.PersistentFlags().StringP("date", "d", "", "Specify date")
	getCmd.PersistentFlags().StringP("user", "u", "", "Specify user")

	rootCmd.AddCommand(getCmd)
}
