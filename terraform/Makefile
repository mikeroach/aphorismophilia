terraform-cmd = TF_IN_AUTOMATION=1 terraform

A13A_VERSION ?= "master-cd6565d-88"
TFVARS_SECRET ?= "../secrets/ephemeral.tfvars"
VARS_SECRET ?= "../secrets/ephemeral.vars"
TFVARS ?= -var="service_aphorismophilia_version=${A13A_VERSION}"

VARS = $(TFVARS) -var-file=$(TFVARS_SECRET)

CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE ?= ../secrets/ephemeral-testing-service-account.json
export CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE
KUBECONFIG ?= ../secrets/ephemeral-kubeconfig
export KUBECONFIG

PROJECT_NAME := ${shell awk -F = '/^gcp_project_shortname/{gsub(/[ |\"]/, ""); print $$2}' $(VARS_SECRET) }
PROJECT_ID := ${shell awk -F = '/^gcp_project_id/{gsub(/[ |\"]/, ""); print $$2}' $(VARS_SECRET) }
DOMAIN := ${shell awk -F = '/^dns_domain/{gsub(/[ |\"]/, ""); print $$2}' $(TFVARS_SECRET) }
HOST := ${shell awk -F = '/^dns_hostname/{gsub(/[ |\"]/, ""); print $$2}' $(TFVARS_SECRET) }
GCP_ZONE := ${shell awk -F = '/^gcp_zone/{gsub(/[ |\"]/, ""); print $$2}' $(VARS_SECRET) }

environment: tf-init tf-fmt tf-validate tf-plan tf-apply
destroy: tf-destroy

k8s-credentials:
	gcloud container clusters get-credentials ${PROJECT_NAME}-k8s-cluster --zone ${GCP_ZONE} --project ${PROJECT_ID}

tf-init:
	#$(terraform-cmd) init -input=false $(VARS)
	if [ -e variables.tf ] ; then $(terraform-cmd) init -input=false $(BACKEND) $(VARS) ;\
		else \
			mv variables.tf variables.tfx ;\
			$(terraform-cmd) init -input=false $(VARS) ;\
			mv variables.tfx variables.tf ;\
	fi

tf-fmt:
	$(terraform-cmd) fmt -check -recursive -diff

tf-validate:
	$(terraform-cmd) validate $(VARS)

tf-plan:
	$(terraform-cmd) plan -input=false $(VARS)

tf-apply:
	$(terraform-cmd) apply -auto-approve -input=false $(VARS)

tf-destroy:
	$(terraform-cmd) destroy -auto-approve -input=false $(VARS)
	rm -f terraform.tfstate terraform.tfstate.backup $(KUBECONFIG)

http-host:
	echo ${HOST}.${DOMAIN}