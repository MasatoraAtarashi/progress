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

// Commits はリポジトリごとのコミット内容とコミット数
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
		commits, err = getProgress(cmd, repository, username, date)
		if commits.Count > 0 {
			repositoryName := strings.Split(repository, "/")
			output += "## " + repositoryName[len(repositoryName)-1] + "(" + strconv.Itoa(commits.Count) + " commits)" + "\n"
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
func getProgress(cmd *cobra.Command, repository string, username string, date string) (commit Commits, err error) {
	startDate, err := time.Parse(layout, date)
	endDate := startDate.AddDate(0, 0, 1)
	cmdArgs := []string{
		"-C", repository, "log",
		"--oneline",
		"--author=" + username,
		"--since=" + startDate.Format(layout),
		"--until=" + endDate.Format(layout),
	}

	// branchオプション
	branch, err := cmd.PersistentFlags().GetString("branch")
	if err != nil {
		return
	}
	if branch != "" {
		cmdArgs = append(cmdArgs, branch)
	} else {
		cmdArgs = append(cmdArgs, "--branches")
	}

	// reverseオプション
	reverse, err := cmd.PersistentFlags().GetBool("reverse")
	if err != nil {
		return
	}
	if reverse {
		cmdArgs = append(cmdArgs, "--reverse")
	}

	// timeオプション
	time, err := cmd.PersistentFlags().GetBool("time")
	if err != nil {
		return
	}
	if time {
		cmdArgs = append(cmdArgs, "--format= - %C(auto)%h%Creset : %s %C(green)(%ad)%Creset")
	} else {
		cmdArgs = append(cmdArgs, "--format= - %C(auto)%h%Creset %s")
	}

	// gitコマンドを実行
	out, err := execGitCmd(cmdArgs)
	commit.Content = string(out)
	commit.Count = len(strings.Split(commit.Content, "\n")) - 1
	return
}

func execGitCmd(cmdArgs []string) (out []byte, err error) {
	cmd := exec.Command(
		"git", cmdArgs...,
	)
	cmd.Stderr = os.Stderr
	out, err = cmd.Output()
	return
}

func init() {
	getCmd.PersistentFlags().StringP("date", "d", "", "Specify date like <2021-04-24>")
	getCmd.PersistentFlags().StringP("user", "u", "", "Specify user")
	getCmd.PersistentFlags().StringP("branch", "b", "", "Specify branch")
	getCmd.PersistentFlags().BoolP("reverse", "r", false, "Reverse order of commits")
	getCmd.PersistentFlags().BoolP("time", "t", false, "Show time")

	rootCmd.AddCommand(getCmd)
}
