# Metrics Bundle Read

consul-debug-read allows users to parse a consul debug bundle quickly to retrieve captured metric values
and display them in a useful output format for quick interpretation

## CLI Usage

| metrics command + flag                                               | return description                                                                                                                                                                                                                                                                                                                                                     | Type          |
  |----------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------|
| `consul-debug-read metrics`/ `consul-debug-read --summary`           | * Datacenter<br/>* Hostname<br/>* Agent Version<br/>* Raft State (Leader/Follower)<br/>* Capture Duration/Interval Start/Stop Timestamps<br/>* Total Number of Metric Scrape Captures                                                                                                                                                                                  | Info          |
| `consul-debug-read metrics --list`                                   | List available metric names to parse with by name.                                                                                                                                                                                                                                                                                                                     | Info          |
| `consul-debug-read metrics --name <metric_name>`                     | Retrieve specific metric timestamped values by name.                                                                                                                                                                                                                                                                                                                   | Value         |
| `consul-debug-read metrics --<metric_group_type>`                    | Retrieves key metric values for Consul metrics pertaining to a particular Consul use-case of diagnosis.<br/>Options include:<br/>* `key-metrics`<br/>* `dataplane`<br/>* `transaction-timing`<br/>* `rate-limiting`<br/>* `leadership-changes`<br/>* `cert-authority`<br/>* `auto-pilot`<br/>* `memory`<br/>* `network`<br/>* `raft-thread-saturation`<br/>* `bolt-db` | Value         |
| `consul-debug-read metrics <value_return_type_flag> --sort-by-value` | Retrieves desired list of metrics by name or group and sorts metric values with highest on top.                                                                                                                                                                                                                                                                        | Value Sorting |



## Metric Grouping


#### *Consul Rate Limiting*

Signs that an agent is being rate-limited or fails to make an RPC request to a 
Consul server can be sudden large changes to the consul.client.rpc metrics(for example, 
greater than 50% deviation from baseline),as well as `consul.client.rpc.exceeded`
and `consul.client.rpc.failed` having a non-zero value.


*Blocking Queries* (Consul Versions <1.10.x)
Rate limit. The blocking query mechanism is reasonably efficient when updates are relatively 
rare (order of tens of seconds to minutes between updates). In cases where a result gets 
updated very fast however - possibly during an outage or incident with a badly behaved 
client - blocking query loops degrade into busy loops that consume excessive client CPU 
and cause high server load. While it's possible to just add a sleep to every iteration of 
the loop, this is not recommended since it causes update delivery to be delayed in the 
happy case, and it can exacerbate the problem since it increases the chance that the index 
has changed on the next request. Clients should instead rate limit the loop so that in the 
happy case they proceed without waiting, but when values start to churn quickly they degrade 
into polling at a reasonable rate (say every 15 seconds). Ideally this is done with an algorithm 
that allows a couple of quick successive deliveries before it starts to limit rate - a token bucket 
with burst of 2 is a simple way to achieve this.

"consul.client.rpc"
"consul.client.rpc.failed"
"consul.client.rpc.exceeded"
"consul.rpc.queries_blocking"
"consul.rpc.rate_limit.exceeded"
"consul.rpc.rate_limit.log_dropped"