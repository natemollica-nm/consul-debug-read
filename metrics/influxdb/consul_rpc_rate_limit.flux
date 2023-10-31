from(bucket: "consul-debug-metrics")
  |> range(start: 2023-10-23T15:03:40Z, stop: 2023-10-23T15:08:30Z)
  |> filter(fn: (r) =>
    r["Name"] == "consul.rpc.rate_limit.exceeded" and
    r["_field"] == "Count"
  )