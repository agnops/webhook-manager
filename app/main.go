package main

import (
	"os"
)

var sslVerification = os.Getenv("sslVerification")

func main() {
	scmProvider := os.Getenv("scmProvider")
	initK8sClientset()
	runner(scmProvider)
}
