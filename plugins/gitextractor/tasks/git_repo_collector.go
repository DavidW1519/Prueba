package tasks

import (
	"strings"
	"errors"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitextractor/parser"
	"github.com/merico-dev/lake/plugins/gitextractor/store"
)

type GitExtractorOptions struct {
	RepoId     string `json:"repoId"`
	Url        string `json:"url"`
	User       string `json:"user"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
	Proxy      string `json:"proxy"`
}

func (o GitExtractorOptions) Valid() error {
	if o.RepoId == "" {
		return errors.New("empty repoId")
	}
	if o.Url == "" {
		return errors.New("empty url")
	}
	url := strings.TrimPrefix(o.Url, "ssh://")
	if !(strings.HasPrefix(o.Url, "http") || strings.HasPrefix(url, "git@") || strings.HasPrefix(o.Url, "/")) {
		return errors.New("wrong url")
	}
	if o.Proxy != "" && !strings.HasPrefix(o.Proxy, "http://") {
		return errors.New("only support http proxy")
	}
	return nil
}

func CollectGitRepo(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	storage := store.NewDatabase(db)
	defer storage.Close()
	op := taskCtx.GetData().(GitExtractorOptions)
	ctx := taskCtx.GetContext()
	p := parser.NewLibGit2(storage)
	var err error
	if strings.HasPrefix(op.Url, "http") {
		err = p.CloneOverHTTP(ctx, op.RepoId, op.Url, op.User, op.Password, op.Proxy)
	} else if url := strings.TrimPrefix(op.Url, "ssh://"); strings.HasPrefix(url, "git@") {
		err = p.CloneOverSSH(ctx, op.RepoId, url, op.PrivateKey, op.Passphrase)
	} else if strings.HasPrefix(op.Url, "/") {
		err = p.LocalRepo(ctx, op.Url, op.RepoId)
	}
	if err != nil {
		return err
	}
	return nil
}
