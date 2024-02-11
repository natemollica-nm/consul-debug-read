SHELL=$(PWD)/shell

run-linter:
	@golangci-lint run ./cmd/cli
	@golangci-lint run ./internal/read

all: clean consul-debug-read split-debug-metrics init-influxdb configure-influxdb telegraf grafana

telemetry: clean init-influxdb configure-influxdb telegraf grafana

consul-debug-read:
	@go install
	@echo "consul-debug-read built and installed => $${GOPATH}/consul-debug-read"

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

clean: stop-grafana stop-telegraf stop-influxdb clean-influxdb clean-grafana

.PHONY:
.SILENT: