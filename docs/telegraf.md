# Consul Offline Telemetry Analysis

[Reference](https://hashicorp.atlassian.net/wiki/spaces/CSE/pages/2317811727/Consul+Offline+Telemetry+Analysis)
Credit: @ranjandas

### Prerequisites

* Telegraf
* InfluxDB

Install and InfluxDB and Telegraf as outlined below in [Getting Started with InfluxDB and Telegraf](#Getting-Started-with-InfluxDB-and-Telegraf).

### Pushing Debug Bundle metrics to InfluxDB using Telegraf

The steps below are meant to simplify the configuration steps required to setup influxDB and telegraf using this repo's root
[Makefile](https://github.com/natemollica-nm/consul-debug-read/blob/main/Makefile) and the following scripts:

* [scripts/init-influxdb.sh](https://github.com/natemollica-nm/consul-debug-read/blob/main/scripts/init-influxdb.sh)
* [scripts/configure-influxdb.sh](https://github.com/natemollica-nm/consul-debug-read/blob/main/scripts/configure-influxdb.sh)
* [scripts/run-telegraf.sh](https://github.com/natemollica-nm/consul-debug-read/blob/main/scripts/run-telegraf.sh)

**1. Setup `consul-debug-read` debug path**

Ensure you've established your terminal session debug bundle path using one of the three methods outlined in
[Working with Debug Bundles](https://github.com/natemollica-nm/consul-debug-read/tree/main#Working-with-Debug-Bundles)

**2. Generate telegraf compatiable metrics.json files**

This step will generate the required debug bundles metrics files in RFC3339 timeformatted files for Telegraf ingestion
at `metrics/telegraf/`

run: `make telegraf-metrics`

**3. Initialize and create influxDB** 

run: `make init-influxdb configure-influxdb`

**4. Ingest metrics with Telegraf**

run: `make telegraf`

The following values are used for this scripted run of telegraf:
* **Username:** consul
* **Password:** hashicorp
* **Initial Organization Name:** hashicorp
* **Initial Bucket Name:** consul-12345 (it is better to suffix the bucket name with the ticket number if you are planning to re-use the InfluxDB instance)

This could take some time depending on the amount of metrics captured during the debug run.

**5. Access InfluxDB UI at http://localhost:8086 to observe ingested debug bundle metrics.**

---
### InfluxDB/Telegraf Teardown

run: `make clean`

This will 
* Stop the Telegraf process
* Stop the InfluxDB process
* Delete the bolt and sqlite databases that influxDB uses to handle the metrics data at:
  * `${HOME}/.influxdbv2/pid` 
  * `${HOME}/.influxdbv2/configs` 
  * `${HOME}/.influxdbv2/influxd.bolt` 
  * `${HOME}/.influxdbv2/influxd.sqlite`

---

## Getting Started with InfluxDB and Telegraf

### Setup InfluxDB

#### Homebrew

```shell
$ brew install influxdb
==> Pouring influxdb--2.7.1.arm64_ventura.bottle.1.tar.gz
....
```

**Auto-start InfluxDB**
To start influxdb now and restart at login:

```shell
$ brew services start influxdb
```

**Manually Start InfluxDB**
Or, if you don't want/need a background service you can just run:

```shell
$ INFLUXD_CONFIG_PATH="/opt/homebrew/etc/influxdb2/config.yml" /opt/homebrew/opt/influxdb/bin/influxd
```

#### Manual Download
```shell
$ curl -LO https://dl.influxdata.com/influxdb/releases/influxdb2-2.1.1-darwin-amd64.tar.gz
$ tar -xzvf  influxdb2-2.1.1-darwin-amd64.tar.gz
$ cd influxdb2-2.1.1-darwin-amd64/
$ ./influxd
...
...
2022-03-20T23:08:17.393403Z	info	Starting	{"log_id": "0_MPFwVW000", "service": "telemetry", "interval": "8h"}
2022-03-20T23:08:17.393706Z	info	Listening	{"log_id": "0_MPFwVW000", "service": "tcp-listener", "transport": "http", "addr": ":8086", "port": 8086}
```

### Setup Telegraf

#### Homebrew

```shell
$ brew install telegraf
==> Pouring telegraf--1.28.2.arm64_ventura.bottle.tar.gz
...
```

**Auto-start Telegraf**
To start telegraf now and restart at login:

```shell
$ brew services start telegraf
```

**Manually Start Telegraf**
Or, if you don't want/need a background service you can just run:

```shell
$ /opt/homebrew/opt/telegraf/bin/telegraf -config metrics/telegraf/telegraf.conf --once --debug
```

#### Install influxdb-cli

```shell
$ brew install influxdb-cli
```

#### Manual Download
```shell
$ curl -LO https://dl.influxdata.com/telegraf/releases/telegraf-1.21.4_darwin_amd64.tar.gz
$ tar -xzvf telegraf-1.21.4_darwin_amd64.tar.gz
$ cd telegraf-1.21.4
```

### DMG Installations

[macOS Intel](https://dl.influxdata.com/telegraf/releases/telegraf-1.28.2_darwin_amd64.dmg?_gl=1*l8nesi*_ga*MTk1MDA5OTg0OC4xNjk2NTMwNDMy*_ga_CNWQ54SDD8*MTY5NzgzMjI2MC42LjEuMTY5NzgzMjU0NC40OC4wLjA.)

[macOS Arm](https://dl.influxdata.com/telegraf/releases/telegraf-1.28.2_darwin_arm64.dmg?_gl=1*l8nesi*_ga*MTk1MDA5OTg0OC4xNjk2NTMwNDMy*_ga_CNWQ54SDD8*MTY5NzgzMjI2MC42LjEuMTY5NzgzMjU0NC40OC4wLjA.)
