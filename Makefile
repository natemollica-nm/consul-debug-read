SHELL=$(PWD)/shell

consul-debug-read:
	@go install

split-debug-metrics:
	@consul-debug-read metrics --telegraf

influxdb: clean init-influxdb configure-influxdb

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

clean-influxdb:
	@echo "removing $${HOME}/.influxdbv2/pid, configs, influxd.bolt, and influxd.sqlite"
	@rm -rf $${HOME}/.influxdbv2/pid
	@rm -rf $${HOME}/.influxdbv2/configs
	@rm -rf $${HOME}/.influxdbv2/influxd.bolt
	@rm -rf $${HOME}/.influxdbv2/influxd.sqlite
	@sleep 8

clean: stop-telegraf stop-influxdb clean-influxdb

.PHONY:
.SILENT: