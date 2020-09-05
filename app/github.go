package main

import (
	"context"
	"flag"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
	"time"
)

var githubClient *github.Client

func InitGitHubClient(accessToken string) context.Context {
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	githubClient = github.NewClient(tc)

	return ctx
}

func CreateGitHubWebhook(accessToken string, webhookUrl string, webhookSecret string, GitOrgProject string, GitRepository string) error {

	sslVrf := "0"
	if strings.ToLower(sslVerification) == "false" {
		sslVrf = "1"
	}

	ctx := InitGitHubClient(accessToken)
	hook := &github.Hook{
		Events: []string{"push", "pull_request", "delete"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          webhookUrl,
			"content_type": "json",
			"secret":       webhookSecret,
			"insecure_ssl": sslVrf,
		},
	}

	if _, _, err := githubClient.Repositories.CreateHook(ctx, GitOrgProject, GitRepository, hook); err != nil {
		// Ignore exists error if the list doesn't return all installed hooks
		if strings.Contains(err.Error(), "Hook already exists on this repository") {
			return nil
		}
		return err
	}

	return nil
}

func GetGitHubAllRepositories(accessToken string) []string {
	ctx := InitGitHubClient(accessToken)
	opt := &github.RepositoryListOptions{Type: "all", Sort: "created", Direction: "desc"} //{Type: "owner", Sort: "updated", Direction: "desc"}

	repos, _, err := githubClient.Repositories.List(ctx, "", opt)
	if err != nil {
		log.Println(err)
	}
	var allRepos []string
	for _, repo := range repos {
		//if !*repo.Fork {
		allRepos = append(allRepos, *repo.Name)
		//}
	}

	time.Sleep(60 * time.Second)

	return allRepos
}

func githubRunner(accessToken string, webhookPayloadUrl string, webhookSecret string, gitOrgProject string)  {
	var repos []string
	repos = GetGitHubAllRepositories(accessToken)

	for _, repo := range repos {
		log.Println(repo)
		CreateGitHubWebhook(accessToken, webhookPayloadUrl, webhookSecret, gitOrgProject, repo)
	}
}