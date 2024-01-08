/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	output                = "./output"
	logger                log.Logger
	ctx                   = context.Background()
	repoMericoLake        = "/Users/houlinwei/Code/go/src/github.com/merico-dev/lake"
	repoMericoLakeWebsite = "/Users/houlinwei/Code/go/src/github.com/merico-dev/website"
	repoId                = "test-repo-id"

	storage        models.Store
	gitRepoCreator *GitRepoCreator

	goGitStorage     models.Store
	goGitRepoCreator *GitRepoCreator
)

func TestMain(m *testing.M) {
	fmt.Println("test main starts")
	logger = logruslog.Global.Nested("git extractor")
	fmt.Println("logger inited")

	var err error
	storage, err = store.NewCsvStore(output + "_libgit2")
	if err != nil {
		panic(err)
	}
	defer storage.Close()
	fmt.Println("git storage inited")
	gitRepoCreator = NewGitRepoCreator(storage, logger)

	goGitStorage, err = store.NewCsvStore(output + "_gogit")
	if err != nil {
		panic(err)
	}
	defer goGitStorage.Close()
	fmt.Println("go git storage inited")
	goGitRepoCreator = NewGitRepoCreator(goGitStorage, logger)

	fmt.Printf("test main run success\n\tlogger: %+v\tstorage: %+v\tgogit storage: %+v\n", logger, storage, goGitStorage)
	m.Run()
}

func TestGitRepo_CountRepoInfo(t *testing.T) {
	repoPath := repoMericoLakeWebsite

	gitRepo, err := gitRepoCreator.LocalRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}
	goGitRepo, err := goGitRepoCreator.LocalGoGitRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}

	tagsCount1, err1 := gitRepo.CountTags(ctx)
	if err1 != nil {
		panic(err1)
	}
	tagsCount2, err2 := goGitRepo.CountTags(ctx)
	if err2 != nil {
		panic(err2)
	}
	t.Logf("[tagsCount] libgit2 result: %d, gogit result: %d", tagsCount1, tagsCount2)
	assert.Equalf(t, tagsCount1, tagsCount2, "unexpected")

	branchesCount1, err1 := gitRepo.CountBranches(ctx)
	if err1 != nil {
		panic(err1)
	}
	branchesCount2, err2 := goGitRepo.CountBranches(ctx)
	if err2 != nil {
		panic(err2)
	}
	t.Logf("[branchesCount] libgit2 result: %d, gogit result: %d", branchesCount1, branchesCount2)
	assert.Equalf(t, branchesCount1, branchesCount2, "unexpected")

	commitCount1, err1 := gitRepo.CountCommits(ctx)
	if err1 != nil {
		panic(err1)
	}
	commitCount2, err2 := goGitRepo.CountCommits(ctx)
	if err2 != nil {
		panic(err2)
	}
	t.Logf("[commitCount] libgit2 result: %d, gogit result: %d", commitCount1, commitCount2)
	assert.Equalf(t, commitCount1, commitCount2, "unexpected")

}

// all testes pass
func TestGitRepo_CollectRepoInfo(t *testing.T) {
	repoPath := repoMericoLake
	gitRepo, err := gitRepoCreator.LocalRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}
	goGitRepo, err := goGitRepoCreator.LocalGoGitRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}

	{
		// finished
		subTaskCtxCollectTags := &testSubTaskContext{}
		if err1 := gitRepo.CollectTags(subTaskCtxCollectTags); err1 != nil {
			panic(err1)
		}
		subTaskCtxCollectTagsWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectTags(subTaskCtxCollectTagsWithGoGit); err2 != nil {
			panic(err2)
		}
		t.Logf("[CollectTags] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectTags, subTaskCtxCollectTagsWithGoGit)
		assert.Equalf(t, subTaskCtxCollectTags.total, subTaskCtxCollectTagsWithGoGit.total, "unexpected")
	}

	{
		// finished
		subTaskCtxCollectBranches := &testSubTaskContext{}
		if err1 := gitRepo.CollectBranches(subTaskCtxCollectBranches); err1 != nil {
			panic(err1)
		}
		subTaskCtxCollectBranchesWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectBranches(subTaskCtxCollectBranchesWithGoGit); err2 != nil {
			panic(err2)
		}
		t.Logf("[CollectBranches] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectBranches, subTaskCtxCollectBranchesWithGoGit)
		assert.Equalf(t, subTaskCtxCollectBranches.total, subTaskCtxCollectBranchesWithGoGit.total, "unexpected")
	}

	{
		// WIP
		subTaskCtxCollectCommits := &testSubTaskContext{}
		if err1 := gitRepo.CollectCommits(subTaskCtxCollectCommits); err1 != nil {
			panic(err1)
		}
		subTaskCtxCCollectCommitsWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectCommits(subTaskCtxCCollectCommitsWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectCommits] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectCommits, subTaskCtxCCollectCommitsWithGoGit)
		fmt.Println(subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total, "unexpected")
		compareTwoStringSlice(b1, b2)
	}

	{
		// TODO CollectDiffLine()
		subTaskCtxCollectDiffLine := &testSubTaskContext{}
		if err1 := gitRepo.CollectDiffLine(subTaskCtxCollectDiffLine); err1 != nil {
			panic(err1)
		}
		subTaskCtxCollectDiffLineWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectCommits(subTaskCtxCollectDiffLineWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectDiffLine] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectDiffLine, subTaskCtxCollectDiffLineWithGoGit)
		fmt.Println(subTaskCtxCollectDiffLine.total, subTaskCtxCollectDiffLineWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectDiffLine.total, subTaskCtxCollectDiffLineWithGoGit.total, "unexpected")
	}
}

func TestGitRepo_CollectCommits(t *testing.T) {
	repoPath := repoMericoLakeWebsite
	gitRepo, err := gitRepoCreator.LocalRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}
	goGitRepo, err := goGitRepoCreator.LocalGoGitRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}

	{
		subTaskCtxCollectCommits := &testSubTaskContext{}
		if err1 := gitRepo.CollectCommits(subTaskCtxCollectCommits); err1 != nil {
			panic(err1)
		}

		subTaskCtxCCollectCommitsWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectCommits(subTaskCtxCCollectCommitsWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectCommits] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectCommits, subTaskCtxCCollectCommitsWithGoGit)
		fmt.Println(subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total, "unexpected")
		compareTwoStringSlice(b1, b2)
	}
}

// compareTwoStringSlice helps to find the difference between two string slices.
func compareTwoStringSlice(b1, b2 []string) {
	for _, b := range b2 {
		var found bool
		for _, bb := range b1 {
			if bb == b {
				found = true
			}
		}
		if !found {
			fmt.Printf("%s from b2, not found in b1\n", b)
		}
	}

	for _, b := range b1 {
		var found bool
		for _, bb := range b2 {
			if bb == b {
				found = true
			}
		}
		if !found {
			fmt.Printf("%s from b1, not found in b2\n", b)
		}
	}
	fmt.Println("compareTwoStringSlice done", len(b1), len(b2))
}
