# Consul Visualization with Telegraf/InfluxDB and Grafana

### Download and Install Grafana

```shell
$ brew update
$ brew install grafana
```

### Start Grafana

```shell
$ brew services start grafana
```

Homebrew configures the following grafana default command:

```shell
$ /opt/homebrew/opt/grafana/bin/grafana server \
    --config /opt/homebrew/etc/grafana/grafana.ini \
    --homepath /opt/homebrew/opt/grafana/share/grafana \
    --packaging\=brew cfg:default.paths.logs\=/opt/homebrew/var/log/grafana \
    cfg:default.paths.data\=/opt/homebrew/var/lib/grafana cfg:default.paths.plugins\=/opt/homebrew/var/lib/grafana/plugins
```

### Using Grafana with InfluxDB

To configure Grafana from the command line without user interaction in the UI, you can use Grafana's built-in command-line tool called `grafana-cli`. You'll be able to perform various configuration tasks, including creating data sources, setting up dashboards, and more. Here's a general procedure for non-interactively configuring Grafana:

**1. Create a Configuration File:**

First, create a configuration file that defines the resources you want to create or configure in Grafana. You can use JSON or YAML formats for this file. Let's say you want to create an InfluxDB data source as an example:

```yaml
apiVersion: 1

datasources:
  - name: consul-debug-metrics
    type: influxdb
    orgId: hashicorp
    access: proxy
    user: consul
    url: http://localhost:8086
    jsonData:
      dbName: consul-debug-grafana
      httpMode: GET
    secureJsonData:
      password: hashicorp
```

Save this file as `grafana-config.yaml`.

**2. Use `grafana-cli` to Apply Configuration:**

Now, you can use the `grafana-cli` tool to apply this configuration. Here's the command to create the InfluxDB data source:

```bash
grafana-cli admin reset-admin-password new-password
grafana-cli admin set-admin-password your-new-password

grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" plugins install grafana-clock-panel
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" plugins install grafana-simple-json-datasource
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" plugins install grafana-worldmap-panel
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" plugins install grafana-piechart-panel
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" plugins install savantly-heatmap-panel
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" plugins install savantly-heatmap-panel

grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" admin reset-admin-password admin
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" admin reset-admin-password your-new-password
grafana-cli --homepath "grafana installation directory" --config "grafana.ini path" admin reset-admin-password new-password

```
Replace the placeholders like `http://localhost:8086`, `mydb`, `username`, and `password` with your specific InfluxDB configuration.

**3. Apply the Configuration:**

Run the following command to apply the configuration:

```bash
grafana-cli admin reset-admin-password your-new-password
grafana-cli admin set-admin-password your-new-password
grafana-cli admin reset-admin-password your-new-password
```

This command will reset the admin password to 'admin', set a new admin password, and then reset it to the password you specified.

**4. Verify the Configuration:**

After applying the configuration, you can verify that the InfluxDB data source has been created by checking the Grafana UI or using `grafana-cli` to export the configuration:

```bash
grafana-cli admin export grafana-config.yaml
```

This command will export the current Grafana configuration, including data sources, into a file named `grafana-config.yaml`.

By following these steps, you can configure Grafana from the command line without any user interaction in the UI. You can similarly define and apply configurations for dashboards, users, and other resources in Grafana using the same approach.