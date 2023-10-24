SHELL=$(PWD)/shell

consul-debug-read:
	@go install

telegraf-metrics:
	@consul-debug-read metrics --telegraf

init-influxdb:
	@scripts/init-influxdb.sh

configure-influxdb:
	@scripts/configure-influxdb.sh

telegraf:
	@scripts/run-telegraf.sh

stop-influxdb:
	@scripts/stop-influxdb.sh

telegraf-token:
	@echo $${TOKEN}

clean-influxdb:
	@rm -rf $${HOME}/.influxdbv2/pid
	@rm -rf $${HOME}/.influxdbv2/configs
	@rm -rf $${HOME}/.influxdbv2/influxd.bolt
	@rm -rf $${HOME}/.influxdbv2/influxd.sqlite

.PHONY:
.SILENT: