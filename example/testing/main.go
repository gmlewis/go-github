// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This embed-interface example is a copy of the "simple" example
// and its purpose is to demonstrate how embedding an interface
// in a struct makes it easy to mock one or more methods.
package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

type RepoReport struct {
	Repo     *github.Repository
	Topics   []string
	Branches []*github.Branch
}

func GenerateReposReport(
	ctx context.Context,
	username string,
	accessToken string,
	ghClient *github.Client,
) ([]*RepoReport, error) {
	repos, _, listReposErr := ghClient.Repositories.List(
		ctx,
		username,
		&github.RepositoryListOptions{
			Visibility: "public",
		},
	)

	if listReposErr != nil {
		return nil, listReposErr
	}

	reports := []*RepoReport{}

	for _, r := range repos[:2] {
		topics, _, topicsErr := ghClient.Repositories.ListAllTopics(
			ctx,
			username,
			*r.Name,
		)

		if topicsErr != nil {
			return nil, topicsErr
		}

		branches, _, branchesErr := ghClient.Repositories.ListBranches(
			ctx,
			username,
			*r.Name,
			&github.BranchListOptions{},
		)

		if branchesErr != nil {
			return nil, branchesErr
		}

		reports = append(reports, &RepoReport{
			Repo:     r,
			Topics:   topics,
			Branches: branches,
		})
	}

	return reports, nil
}

func main() {
	var username string
	var accessToken string

	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	fmt.Print("Enter GitHub access token (you can create one at https://github.com/settings/tokens): ")
	fmt.Scanf("%s", &accessToken)

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	fmt.Println("Generating report...")
	repoReports, err := GenerateReposReport(
		ctx,
		username,
		accessToken,
		client,
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\nREPOSITORIES REPORT:")
	for _, repoReport := range repoReports {
		fmt.Println("Repo: ", *repoReport.Repo.Name)
		fmt.Println("Topics: ", repoReport.Topics)
		fmt.Printf("Branches:")
		for _, b := range repoReport.Branches {
			fmt.Printf(" " + *b.Name)
		}
		fmt.Printf("\n")
	}
}
