package main

import (
	"fmt"
	"os"
)

func detectCloudShell() string {
	if os.Getenv("AWS_EXECUTION_ENV") == "CloudShell" {
		return "AWS Cloud Shell"
	}
	if os.Getenv("CLOUD_SHELL") == "true" {
		return "GCP Cloud Shell"
	}
	if os.Getenv("ACC_TERM") != "" || os.Getenv("AZUREPS_HOST_ENVIRONMENT") != "" {
		return "Azure Cloud Shell"
	}
	return "local environment"
}

func main() {
	fmt.Println("Running in", detectCloudShell())
}
