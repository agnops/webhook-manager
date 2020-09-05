package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var clientset *kubernetes.Clientset
var secretsClient coreV1Types.SecretInterface

var namespace = os.Getenv("NAMESPACE")

func initK8sClientset() {
	var err error
	var config *rest.Config

	// use the current context in kubeconfig
	config, err = rest.InClusterConfig()

	if err != nil {
		if  strings.Contains(err.Error(), "unable to load in-cluster configuration") {
			var kubeconfig *string
			if home := homeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		} else {
			log.Println(err.Error())
		}
	}
	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	failOnError(err, "Failed on clientset init")
	secretsClient = clientset.CoreV1().Secrets(namespace)
}

func getRegisteredOrgUser(scmProvider string) []apiv1.Secret {
	secretList, err := secretsClient.List(context.TODO(), metav1.ListOptions{LabelSelector: "scmProvider=" + scmProvider})
	if err != nil {
		log.Println(err.Error())
	}
	return secretList.Items
}

func getWebhookSecret(secretName string) string {
	secret, err := secretsClient.Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	if secret == nil || secret.GetName() != secretName {
		webhookSecret, _ := randomHex(20)
		createWebhookSecret(secretName, webhookSecret)
		return webhookSecret
	}
	return string(secret.Data["WebhookSecret"])
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func createWebhookSecret(secretName string, webhookSecret string) {
	secretSpec := apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				"AgnOps":      "WebhookSecret",
			},
		},
		Data: map[string][]byte{
			"WebhookSecret":   []byte(webhookSecret),
		},
		Type: "Opaque",
	}

	secretName = secretSpec.ObjectMeta.Name

	_, err := secretsClient.Create(context.TODO(), &secretSpec, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Created secret %s\n", secretName)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}