#!/bin/bash

gcloud functions deploy set_web_hook |
gcloud functions deploy handle_message |
gcloud functions deploy pool_jira_issues |
gcloud functions deploy expire_notifications |
gcloud functions deploy check_issues_duedate
