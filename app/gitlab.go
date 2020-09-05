package main

import (
	"flag"
	"github.com/xanzy/go-gitlab"
	"log"
	"strings"
)

var gitlabClient *gitlab.Client

func InitGitLabClient(accessToken string) {
	flag.Parse()
	gitlabClient, _ = gitlab.NewOAuthClient(accessToken)
}

func CheckIfWebhookExist(projectID int, webhookUrl string) bool {

	opt := gitlab.ListProjectHooksOptions{}

	for {
		ps, resp, err := gitlabClient.Projects.ListProjectHooks(projectID, &opt)
		if err != nil { log.Fatal(err) }

		for _, p := range ps { if p.URL == webhookUrl { return true } }

		if resp.CurrentPage >= resp.TotalPages { break }
		opt.Page = resp.NextPage
	}

	return false
}

func CreateGitLabWebhook(accessToken string, webhookUrl string, projectID int, webhookSecret string) error {
	InitGitLabClient(accessToken)

	sslVrf := true
	if strings.ToLower(sslVerification) == "false" {
		sslVrf = false
	}

	if !CheckIfWebhookExist(projectID, webhookUrl) {
		projectHookOptions := gitlab.AddProjectHookOptions{
			URL:                      gitlab.String(webhookUrl),
			ConfidentialNoteEvents:   gitlab.Bool(false),
			PushEvents:               gitlab.Bool(true),
			IssuesEvents:             gitlab.Bool(false),
			ConfidentialIssuesEvents: gitlab.Bool(false),
			MergeRequestsEvents:      gitlab.Bool(false),
			TagPushEvents:            gitlab.Bool(false),
			NoteEvents:               gitlab.Bool(false),
			JobEvents:                gitlab.Bool(false),
			PipelineEvents:           gitlab.Bool(false),
			WikiPageEvents:           gitlab.Bool(false),
			EnableSSLVerification:    gitlab.Bool(sslVrf),
			Token:                    gitlab.String(webhookSecret),
		}
		_, _, err := gitlabClient.Projects.AddProjectHook(projectID, &projectHookOptions)
		if err != nil { return err }
	}
	return nil
}

func GetGitLabAllRepositories(accessToken string) []int {
	InitGitLabClient(accessToken)

	opt := &gitlab.ListProjectsOptions{
		Owned: gitlab.Bool(true),
	}

	var allRepos []int

	for {
		// Get the first page with projects.
		ps, resp, err := gitlabClient.Projects.ListProjects(opt)
		if err != nil {
			log.Fatal(err)
		}

		// List all the projects we've found so far.
		for _, p := range ps {
			//fmt.Printf("Found project: %s\n", p.Name)
			//allRepos = append(allRepos, p.Name)
			allRepos = append(allRepos, p.ID)
		}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return allRepos
}

func gitlabRunner(accessToken string, webhookPayloadUrl string, webhookSecret string)  {
	var repos []int
	repos = GetGitLabAllRepositories(accessToken)

	for _, repo := range repos {
		log.Println(repo)
		_ = CreateGitLabWebhook(accessToken, webhookPayloadUrl, repo, webhookSecret)
	}
}