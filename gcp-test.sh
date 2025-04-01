#!/bin/bash

go build -o go-onboarding
gcloud cloud-shell scp localhost:./go-onboarding cloudshell:~/go-onboarding
gcloud cloud-shell ssh --command=./go-onboarding --authorize-session