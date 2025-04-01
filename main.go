package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// Azure
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling"
	
	// GCP
	"google.golang.org/api/cloudbilling/v1"
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
	switch detectCloudShell() {
	case "AWS Cloud Shell":
		fmt.Println("Running in AWS Cloud Shell")
		fmt.Println("Not implemented yet")
		os.Exit(1)
	case "GCP Cloud Shell":
		fmt.Println("Running in GCP Cloud Shell")
		doGoogleOnboarding()
		os.Exit(0)
	case "Azure Cloud Shell":
		fmt.Println("Running in Azure Cloud Shell")
		doAzureOnboarding()
		os.Exit(0)
	default:
		fmt.Println("Running in local environment")
		fmt.Println("Not supported")
		os.Exit(1)
	}
}

func doAzureOnboarding() {
	fmt.Println("Azure Cloud Shell onboarding started")

	// Understanding the time taken for onboarding process
	startTime := time.Now()

	// Create a credential using DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain credentials: %v", err)
	}
	fmt.Printf("Azure Cloud Shell onboarding completed: %v\n", cred)

	// Create billing accounts client
	billingClient, err := armbilling.NewAccountsClient(cred, nil)
	if err != nil {
		log.Fatalf("failed to create billing client: %v", err)
	}
	fmt.Printf("Billing client created: %v\n", billingClient)

	// List billing accounts
	pager := billingClient.NewListPager(nil)
	ctx := context.Background()

	fmt.Sprintf("DEBUG: %v", pager)

	billingAccounts := []*armbilling.Account{} // List of billing accounts
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to get billing accounts: %v", err)
		}

		// Process each billing account
		for _, account := range page.Value {
			if account.ID != nil && account.Name != nil {
				fmt.Printf("Billing Account ID: %s\n", *account.ID)
				fmt.Printf("Billing Account Name: %s\n", *account.Name)
				billingAccounts = append(billingAccounts, account)
			}
		}
	}

	var billingAccount *armbilling.Account
	// Check how many billing accounts were found
	if len(billingAccounts) == 0 {
		fmt.Println("No billing accounts found. You potentially are not signed in as the Billing Account Admin")
		os.Exit(1)
	} else if len(billingAccounts) > 1 {
		fmt.Printf("Multiple billing accounts found: %d\n", len(billingAccounts))
		fmt.Println("Selecting a billing account has not been implemented yet")
		os.Exit(1)
	} else {
		billingAccount = billingAccounts[0]
	}

	fmt.Printf("Selected billing account: %s\n", *billingAccount.Name)

	// Log total time taken for onboarding process
	log.Printf("Azure Cloud Shell onboarding completed in %v", time.Since(startTime))

	// Next steps
	// Attempt to get 1yr cost estimate -- need this for decision on API or Export Based	
}

func doGoogleOnboarding() {
	fmt.Println("GCP Cloud Shell onboarding started")

	// Understanding the time taken for onboarding process
	startTime := time.Now()

	// Set up the Cloud Billing API client using the default credentials
	ctx := context.Background()
	log.Println("Creating Cloud Billing API client...")
	billingService, err := cloudbilling.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to create Cloud Billing service: %v", err)
		return
	}

	// Log that the client has been successfully created
	log.Println("Cloud Billing API client created successfully.")

	// List billing accounts
	billingAccountsService := cloudbilling.NewBillingAccountsService(billingService)
	log.Println("Fetching billing accounts...")
	billingAccountsListCall := billingAccountsService.List()
	billingAccountsListCall = billingAccountsListCall.PageSize(10) // Adjust page size as needed

	billingAccounts, err := billingAccountsListCall.Do()
	if err != nil {
		log.Fatalf("Unable to list billing accounts: %v", err)
		return
	}

	// Log the billing accounts found
	log.Printf("Found %d billing accounts", len(billingAccounts.BillingAccounts))

	// Print out the Billing Account IDs
	for _, billingAccount := range billingAccounts.BillingAccounts {
		fmt.Printf("Billing Account ID: %s\n", billingAccount.Name)
	}

	// Log total time taken for onboarding process
	log.Printf("GCP Cloud Shell onboarding completed in %v", time.Since(startTime))
}