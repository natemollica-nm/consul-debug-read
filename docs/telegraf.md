# Consul Offline Telemetry Analysis

Credit: @ranjandas
Reference: https://hashicorp.atlassian.net/wiki/spaces/CSE/pages/2317811727/Consul+Offline+Telemetry+Analysis


### Prerequisites

Install and configure InfluxDB and Telegraf as outlined in the above [reference](https://hashicorp.atlassian.net/wiki/spaces/CSE/pages/2317811727/Consul+Offline+Telemetry+Analysis)
or as shown below in [Setup InfluxDB](#Setup-InfluxDB) and [Setup Telegraf](#Setup-Telegraf).

#### Generate telegraf compatiable metrics.json files
this will generate the debug bundles metrics files in RFC3339 timeformatted files for Telegraf ingestion
at `metrics/telegraf/`

run: `make telegraf-metrics`

#### Initialize and create influxDB 

run: `make init-influxdb configure-influxdb`

#### Ingest metrics with Telegraf

run: `make telegraf`

This could take some time depending on the amount of metrics captured during the debug run.

## Setup InfluxDB

### Homebrew

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

### Manual Download
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

The following values are used for this example:
* **Username:** consul
* **Password:** hashicorp
* **Initial Organization Name:** hashicorp
* **Initial Bucket Name:** consul-12345 (it is better to suffix the bucket name with the ticket number if you are planning to re-use the InfluxDB instance)

## Setup Telegraf

### Homebrew

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

### Manual Download
```shell
$ curl -LO https://dl.influxdata.com/telegraf/releases/telegraf-1.21.4_darwin_amd64.tar.gz
$ tar -xzvf telegraf-1.21.4_darwin_amd64.tar.gz
$ cd telegraf-1.21.4
```

## DMG Installations

[macOS Intel](https://dl.influxdata.com/telegraf/releases/telegraf-1.28.2_darwin_amd64.dmg?_gl=1*l8nesi*_ga*MTk1MDA5OTg0OC4xNjk2NTMwNDMy*_ga_CNWQ54SDD8*MTY5NzgzMjI2MC42LjEuMTY5NzgzMjU0NC40OC4wLjA.)

[macOS Arm](https://dl.influxdata.com/telegraf/releases/telegraf-1.28.2_darwin_arm64.dmg?_gl=1*l8nesi*_ga*MTk1MDA5OTg0OC4xNjk2NTMwNDMy*_ga_CNWQ54SDD8*MTY5NzgzMjI2MC42LjEuMTY5NzgzMjU0NC40OC4wLjA.)
