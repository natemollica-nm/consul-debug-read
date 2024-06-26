SHELL=$(PWD)/shell
DOMAIN_NAME=nathan-mollica.sbx.hashidemos.io
SUBDOMAIN=consul-support-demo
PUBLIC_IP=$(shell curl -s ifconfig.me)

# /////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Doormat /////////////////////////////////// #
##@ Doormat
.PHONY: doormat-refresh
doormat-refresh: ## Refresh local .env file doormat-credentials
	@scripts/doormat-update-creds.sh

# /////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Linter /////////////////////////////////// #
##@ Linter
run-linter: ## Run golang linter
	@golangci-lint run ./cmd/cli
	@golangci-lint run ./internal/read

# ///////////////////////////////////////////////////////////////////// #
# //////////////////////////// Demo /////////////////////////////////// #
##@ Demo
packer-ami: ## Build demo bundle webserver AMI
	@packer init .offsite/packer/webserver.pkr.hcl
	@packer build .offsite/packer/webserver.pkr.hcl

# ////////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Terraform /////////////////////////////////// #
##@ Terraform
.PHONY: plan
plan: ## Run terraform plan in .offsite dir
	@terraform -chdir=.offsite init
	@terraform -chdir=.offsite plan \
		-var aws_region="$$AWS_REGION" \
		-var subdomain="$(SUBDOMAIN)" \
		-var ssh_key_name="$$AWS_EC2_SSH_KEY_NAME" \
		-var local_public_cidr="$(PUBLIC_IP)/32"

.PHONY: refresh
refresh: ## Run terraform refresh in .offsite dir
	@terraform -chdir=.offsite refresh \
		-var aws_region="$$AWS_REGION" \
		-var subdomain="$(SUBDOMAIN)" \
		-var ssh_key_name="$$AWS_EC2_SSH_KEY_NAME" \
		-var local_public_cidr="$(PUBLIC_IP)/32"

.PHONY: apply
apply: ## Run terraform auto-approved apply in .offsite dir
	@terraform -chdir=.offsite apply \
		-var aws_region="$$AWS_REGION" \
		-var subdomain="$(SUBDOMAIN)" \
		-var ssh_key_name="$$AWS_EC2_SSH_KEY_NAME" \
		-var local_public_cidr="$(PUBLIC_IP)/32" \
		-auto-approve

.PHONY: destroy
destroy: ## Run terraform destroy on aws resources
	@terraform -chdir=.offsite destroy \
		-var aws_region="$$AWS_REGION" \
		-var subdomain="$(SUBDOMAIN)" \
		-var ssh_key_name="$$AWS_EC2_SSH_KEY_NAME" \
		-var local_public_cidr="$(PUBLIC_IP)/32" \
		-auto-approve

.PHONY: destroy-instance
destroy-instance: ## Target destroy bastion host only from aws
	@terraform -chdir=.offsite destroy \
		-var aws_region="$$AWS_REGION" \
		-var subdomain="$(SUBDOMAIN)" \
		-var ssh_key_name="$$AWS_EC2_SSH_KEY_NAME" \
		-var local_public_cidr="$(PUBLIC_IP)/32" \
		-target=aws_instance.web \
		-auto-approve
# /////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Build /////////////////////////////////// #
##@ Build
consul-debug-read: ## Run go install for consul-debug-read cli
	@go install
	@echo "consul-debug-read built and installed => $${GOPATH}/consul-debug-read"

# ///////////////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Telegraf/Grafana /////////////////////////////////// #
##@ Telegraf/Grafana
all: clean consul-debug-read split-debug-metrics init-influxdb configure-influxdb telegraf grafana ## Reset telemetry tools/db, run cli tool metrics formatter, and start new telemetry tooling

telemetry: clean init-influxdb configure-influxdb telegraf grafana

split-debug-metrics:
	@consul-debug-read metrics --telegraf

init-influxdb:
	@scripts/init-influxdb.sh

configure-influxdb:
	@scripts/configure-influxdb.sh

telegraf:
	@scripts/stop-telegraf.sh
	@scripts/run-telegraf.sh

stop-influxdb:
	@scripts/stop-influxdb.sh

stop-telegraf:
	@scripts/stop-telegraf.sh

telegraf-token:
	@echo $${TOKEN}

grafana:
	@scripts/init-grafana.sh

stop-grafana:
	@brew services stop grafana

clean-grafana:
	@echo "removing grafana db and conf data"
	@rm /opt/homebrew/share/grafana >/dev/null 2>&1 || true
	@rm -rf /opt/homebrew/opt/grafana/share/grafana/data/grafana.db

nuke-grafana:
	@echo "nuking grafana brew installation...."
	@brew services stop grafana >/dev/null 2>&1 || true
	@brew unlink grafana >/dev/null 2>&1 || true
	@brew uninstall grafana --zap >/dev/null 2>&1 || true
	@rm -rf /opt/homebrew/etc/grafana/
	@brew cleanup >/dev/null 2>&1 || true
	@echo "grafana nuked"

clean-influxdb:
	@echo "removing influxdb conf, bolt, and sqlite data"
	@rm -rf $${HOME}/.influxdbv2/pid
	@rm -rf $${HOME}/.influxdbv2/configs
	@rm -rf $${HOME}/.influxdbv2/engine
	@rm -rf $${HOME}/.influxdbv2/influxd.bolt
	@rm -rf $${HOME}/.influxdbv2/influxd.sqlite
	@sleep 3

# //////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Cleanup /////////////////////////////////// #
##@ Cleanup
clean: stop-grafana stop-telegraf stop-influxdb clean-influxdb clean-grafana ## Stop telemetry tools and clean influxdb

# /////////////////////////////////////////////////////////////////////////// #
# //////////////////////////// Help Goals /////////////////////////////////// #
.DEFAULT_GOAL := help
##@ Help
# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

%:
	@: