// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-github/v34/github"
)

type fakeRepoSvc struct {
	// embedding the interface for testing purposes
	// makes our lifes easier
	github.RepositoriesServiceInterface

	ListFn func(
		ctx context.Context,
		user string,
		opts *github.RepositoryListOptions,
	) (
		[]*github.Repository,
		*github.Response,
		error,
	)

	ListAllTopicsFn func(
		ctx context.Context,
		owner,
		repo string,
	) (
		[]string,
		*github.Response,
		error,
	)

	ListBranchesFn func(
		ctx context.Context,
		owner string,
		repo string,
		opts *github.BranchListOptions,
	) (
		[]*github.Branch,
		*github.Response,
		error,
	)
}

func (svc *fakeRepoSvc) List(
	ctx context.Context,
	user string,
	opts *github.RepositoryListOptions,
) (
	[]*github.Repository,
	*github.Response,
	error,
) {
	return svc.ListFn(
		ctx,
		user,
		opts,
	)
}

func (svc *fakeRepoSvc) ListAllTopics(
	ctx context.Context,
	owner,
	repo string,
) (
	[]string,
	*github.Response,
	error,
) {
	return svc.ListAllTopicsFn(
		ctx,
		owner,
		repo,
	)
}

func (svc *fakeRepoSvc) ListBranches(
	ctx context.Context,
	owner string,
	repo string,
	opts *github.BranchListOptions,
) (
	[]*github.Branch,
	*github.Response,
	error,
) {
	return svc.ListBranchesFn(
		ctx,
		owner,
		repo,
		opts,
	)
}

func TestCreateReposReport(t *testing.T) {
	tt := []struct {
		name string

		// inputs
		ctx context.Context

		// overrides
		ListFn func(
			ctx context.Context,
			user string,
			opts *github.RepositoryListOptions,
		) (
			[]*github.Repository,
			*github.Response,
			error,
		)

		ListAllTopicsFn func(
			ctx context.Context,
			owner,
			repo string,
		) (
			[]string,
			*github.Response,
			error,
		)

		ListBranchesFn func(
			ctx context.Context,
			owner string,
			repo string,
			opts *github.BranchListOptions,
		) (
			[]*github.Branch,
			*github.Response,
			error,
		)

		// outputs
		wantReport []*RepoReport
		wantErr    error
	}{
		{
			name: "HappyPath",
			ctx:  context.Background(),
			ListFn: func(
				ctx context.Context,
				user string,
				opts *github.RepositoryListOptions,
			) (
				[]*github.Repository,
				*github.Response,
				error,
			) {
				return []*github.Repository{
					{
						Name: github.String("myrepo1"),
					},
					{
						Name: github.String("myrepo2"),
					},
				}, nil, nil
			},
			ListAllTopicsFn: func(
				ctx context.Context,
				owner,
				repo string,
			) (
				[]string,
				*github.Response,
				error,
			) {
				return []string{"topic1", "topic2"}, nil, nil
			},
			ListBranchesFn: func(
				ctx context.Context,
				owner,
				repo string,
				opts *github.BranchListOptions,
			) (
				[]*github.Branch,
				*github.Response,
				error,
			) {
				return []*github.Branch{
					{
						Name: github.String("branch1"),
					},
				}, nil, nil
			},
			wantReport: []*RepoReport{
				{
					Repo: &github.Repository{
						Name: github.String("myrepo1"),
					},
					Topics: []string{"topic1", "topic2"},
					Branches: []*github.Branch{
						{
							Name: github.String("branch1"),
						},
					},
				},
				{
					Repo: &github.Repository{
						Name: github.String("myrepo2"),
					},
					Topics: []string{"topic1", "topic2"},
					Branches: []*github.Branch{
						{
							Name: github.String("branch1"),
						},
					},
				},
			},
		},
		{
			name: "ErrorRepoListing",
			ctx:  context.Background(),
			ListFn: func(
				ctx context.Context,
				user string,
				opts *github.RepositoryListOptions,
			) (
				[]*github.Repository,
				*github.Response,
				error,
			) {
				return nil, nil, errors.New("some error")
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "ErrTopicsListing",
			ctx:  context.Background(),
			ListFn: func(
				ctx context.Context,
				user string,
				opts *github.RepositoryListOptions,
			) (
				[]*github.Repository,
				*github.Response,
				error,
			) {
				return []*github.Repository{
					{
						Name: github.String("myrepo1"),
					},
					{
						Name: github.String("myrepo2"),
					},
				}, nil, nil
			},
			ListAllTopicsFn: func(
				ctx context.Context,
				owner,
				repo string,
			) (
				[]string,
				*github.Response,
				error,
			) {
				return []string{}, nil, errors.New("some error")
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "ErrBranchListing",
			ctx:  context.Background(),
			ListFn: func(
				ctx context.Context,
				user string,
				opts *github.RepositoryListOptions,
			) (
				[]*github.Repository,
				*github.Response,
				error,
			) {
				return []*github.Repository{
					{
						Name: github.String("myrepo1"),
					},
					{
						Name: github.String("myrepo2"),
					},
				}, nil, nil
			},
			ListAllTopicsFn: func(
				ctx context.Context,
				owner,
				repo string,
			) (
				[]string,
				*github.Response,
				error,
			) {
				return []string{"topic1", "topic2"}, nil, nil
			},
			ListBranchesFn: func(
				ctx context.Context,
				owner,
				repo string,
				opts *github.BranchListOptions,
			) (
				[]*github.Branch,
				*github.Response,
				error,
			) {
				return []*github.Branch{}, nil, errors.New("some error")
			},
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tt {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// the http client doesn't reallly matter
			// as we are completely mocking the real method calls
			ghC := github.NewClient(nil)

			// override Resitories with our fake implementations
			ghC.Repositories = &fakeRepoSvc{
				ListFn:          tt.ListFn,
				ListAllTopicsFn: tt.ListAllTopicsFn,
				ListBranchesFn:  tt.ListBranchesFn,
			}

			report, repErr := GenerateReposReport(
				tt.ctx,
				"myusername",
				"myaccesstoken",
				ghC,
			)

			for i, r := range report {
				if !reflect.DeepEqual(r, tt.wantReport[i]) {
					t.Errorf("got = %#v, want = %#v", r, tt.wantReport[i])
				}
			}

			if tt.wantErr == nil {
				if repErr != nil {
					t.Errorf("got = %#v, want = %#v", repErr, tt.wantErr)
				}
			} else {
				if errors.Is(repErr, tt.wantErr) {
					t.Errorf("got = %#v, want = %#v", repErr, tt.wantErr)
				}
			}
		})
	}
}
