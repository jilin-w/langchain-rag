package git_service

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/twitchyliquid64/golang-asm/objabi"
)

// clone仓库
func CloneRepo(repoUrl string) (filePath string, err error) {
	fileName := strings.Split(strings.Split(repoUrl, "/")[len(strings.Split(repoUrl, "/"))-1], ".")[0]
	gitDir := objabi.WorkingDir() + "/tmp/repo/" + fileName
	_, err = os.Stat(gitDir)
	if err == nil || !os.IsNotExist(err) {
		os.RemoveAll(gitDir)
	}
	os.MkdirAll(gitDir, os.ModePerm)
	//拉去
	repo, err := git.PlainClone(gitDir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", err
	}
	_, err = repo.Worktree()
	if err != nil {
		return
	}
	return gitDir, nil
}
