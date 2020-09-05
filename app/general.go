package main

import (
	v1 "k8s.io/api/core/v1"
	"log"
	"os"
	"strings"
	"time"
)

func GetWebhookPayloadUrlAndSecret(scmProvider string) (string, string) {
	webhookPayloadUrl := ""
	webhookSecret := ""

	if (os.Getenv("WEBHOOK_PAYLOAD_URL") != webhookPayloadUrl) && (len(os.Getenv("WEBHOOK_PAYLOAD_URL")) != 0) {
		webhookPayloadUrl = os.Getenv("WEBHOOK_PAYLOAD_URL")
	}
	if (os.Getenv("WEBHOOK_SECRET") != webhookSecret) && (len(os.Getenv("WEBHOOK_SECRET")) != 0) {
		webhookSecret = os.Getenv("WEBHOOK_SECRET")
	} else {
		var helmRelease = os.Getenv("HELM_RELEASE")
		secretName := strings.ToLower(helmRelease) + "-" + scmProvider + "-agnops-webhook-secret"
		webhookSecret = getWebhookSecret(secretName)
	}
	return webhookPayloadUrl, webhookSecret
}

func runner(scmProvider string)  {
	webhookPayloadUrl, webhookSecret := GetWebhookPayloadUrlAndSecret(scmProvider)

	log.Printf("webhookPayloadUrl: %s webhookSecret: %s\n", webhookPayloadUrl, webhookSecret)

	var keys []v1.Secret

	for {
		keys = getRegisteredOrgUser(scmProvider)
		if len(keys) > 0 {
			log.Printf("\nFound %d %s accounts:\n", len(keys), scmProvider)
			for _, key := range keys {
				log.Println(string(key.Data["OrgUserName"]))
				accessToken := string(key.Data["OAuth2Token"])

				switch scmProvider {
				case "github":
					gitOrgProject := string(key.Data["OrgUserName"])
					githubRunner(accessToken, webhookPayloadUrl, webhookSecret, gitOrgProject)
				case "gitlab":
					gitlabRunner(accessToken, webhookPayloadUrl, webhookSecret)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}