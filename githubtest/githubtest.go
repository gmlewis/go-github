// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package githubtest provides an easy way to test code that
// uses the github package.
package githubtest

import (
	"github.com/google/go-github/v35/github"
)

// Option represents an option that can be passed in to New
// to modify the behavior of the returned test client.
type Option func(c *github.Client)

// New returns a new *github.Client with optional (fake)
// services (that you provide) replacing existing ones.
func New(opts ...Option) *github.Client {
	c := github.NewClient(nil)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithOrganizations(fake *github.OrganizationsService) Option {
	return func(c *github.Client) {
		c.Organizations = fake
	}
}

/*
	c.Actions = (*ActionsService)(&c.common)
	c.Activity = (*ActivityService)(&c.common)
	c.Admin = (*AdminService)(&c.common)
	c.Apps = (*AppsService)(&c.common)
	c.Authorizations = (*AuthorizationsService)(&c.common)
	c.Billing = (*BillingService)(&c.common)
	c.Checks = (*ChecksService)(&c.common)
	c.CodeScanning = (*CodeScanningService)(&c.common)
	c.Enterprise = (*EnterpriseService)(&c.common)
	c.Gists = (*GistsService)(&c.common)
	c.Git = (*GitService)(&c.common)
	c.Gitignores = (*GitignoresService)(&c.common)
	c.Interactions = (*InteractionsService)(&c.common)
	c.IssueImport = (*IssueImportService)(&c.common)
	c.Issues = (*IssuesService)(&c.common)
	c.Licenses = (*LicensesService)(&c.common)
	c.Marketplace = &MarketplaceService{client: c}
	c.Migrations = (*MigrationService)(&c.common)
	c.Organizations = (*OrganizationsService)(&c.common)
	c.Projects = (*ProjectsService)(&c.common)
	c.PullRequests = (*PullRequestsService)(&c.common)
	c.Reactions = (*ReactionsService)(&c.common)
	c.Repositories = (*RepositoriesService)(&c.common)
	c.Search = (*SearchService)(&c.common)
	c.Teams = (*TeamsService)(&c.common)
	c.Users = (*UsersService)(&c.common)

var _ ActionsServiceInterface = &github.ActionsService{}
var _ ActivityServiceInterface = &github.ActivityService{}
var _ AdminServiceInterface = &github.AdminService{}
var _ AppsServiceInterface = &github.AppsService{}
var _ AuthorizationsServiceInterface = &github.AuthorizationsService{}
var _ BillingServiceInterface = &github.BillingService{}
var _ ChecksServiceInterface = &github.ChecksService{}
var _ CodeScanningServiceInterface = &github.CodeScanningService{}
var _ EnterpriseServiceInterface = &github.EnterpriseService{}
var _ GistsServiceInterface = &github.GistsService{}
var _ GitServiceInterface = &github.GitService{}
var _ GitignoresServiceInterface = &github.GitignoresService{}
var _ InteractionsServiceInterface = &github.InteractionsService{}
var _ IssueImportServiceInterface = &github.IssueImportService{}
var _ IssuesServiceInterface = &github.IssuesService{}
var _ LicensesServiceInterface = &github.LicensesService{}
var _ MarketplaceServiceInterface = &github.MarketplaceService{}
var _ MigrationServiceInterface = &github.MigrationService{}
var _ OrganizationsServiceInterface = &github.OrganizationsService{}
var _ ProjectsServiceInterface = &github.ProjectsService{}
var _ PullRequestsServiceInterface = &github.PullRequestsService{}
var _ ReactionsServiceInterface = &github.ReactionsService{}
var _ RepositoriesServiceInterface = &github.RepositoriesService{}
var _ SearchServiceInterface = &github.SearchService{}
var _ TeamsServiceInterface = &github.TeamsService{}
var _ UsersServiceInterface = &github.UsersService{}
*/
