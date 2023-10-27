SHELL=$(PWD)/shell

consul-debug-read:
	@go install

split-debug-metrics:
	@consul-debug-read metrics --telegraf

all: clean init-influxdb configure-influxdb telegraf grafana

init-influxdb:
	@scripts/init-influxdb.sh

configure-influxdb:
	@scripts/configure-influxdb.sh

telegraf:
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
	@echo "removing grafana data and conf data"
	@rm -rf /opt/homebrew/opt/grafana/share/grafana/data/grafana.db
	@rm -rf /opt/homebrew/opt/grafana/share/grafana/conf/dashboards/
	@rm -rf /opt/homebrew/opt/grafana/share/grafana/conf/datasources/
	@rm -rf /opt/homebrew/opt/grafana/share/grafana/conf/plugins/
	@sleep 3

clean-influxdb:
	@echo "removing $${HOME}/.influxdbv2/pid, configs, influxd.bolt, and influxd.sqlite"
	@rm -rf $${HOME}/.influxdbv2/pid
	@rm -rf $${HOME}/.influxdbv2/configs
	@rm -rf $${HOME}/.influxdbv2/influxd.bolt
	@rm -rf $${HOME}/.influxdbv2/influxd.sqlite
	@sleep 8

clean: stop-grafana stop-telegraf stop-influxdb clean-influxdb clean-grafana

.PHONY:
.SILENT: