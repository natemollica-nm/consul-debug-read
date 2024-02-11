## Consul Debug Bundle CLI Tool

The consul-debug-read is a command-line-interface tool developed to help
HashiCorp Support Engineers better interpret and make use of a customer
provided `consul debug` command output bundle.

This tool does make the assumption that the engineer using the tool 
has a fundamental understanding of Consul as a service-discovery,
service-mesh, and kev-value store.

## Bundle Contents

The `consul debug` command was implemented with the following purpose in mind:

> _Monitor a Consul agent for a specified period of time, recording
information about the agent, cluster, and environment to an archive
written to the specified path._

Keep in mind the bundle is agent specific, meaning the output capture from
the command is highly dependent upon whether the command was ran
from a Consul client agent or Consul server agent.

In short, nothing in regard to the Consul cluster's Raft status can be expected
to be captured by a debug bundle from a Consul client agent. In the same 
sense, nothing in regard to specific service and service-mesh proxy information
can be expected to be gathered from a Consul server agent. This is important
to understand when requested a bundle from customers.

### Debug Capture Content Breakdown



```shell
├── 2023-11-16T15-12-30Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-13-00Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-13-30Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-14-00Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-14-30Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-15-00Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-15-30Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-16-00Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-16-30Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-17-00Z
│         ├── goroutine.prof
│         └── heap.prof
├── 2023-11-16T15-17-30Z
│         ├── goroutine.prof
│         └── heap.prof
├── agent.json
├── consul.log
├── host.json
├── index.json
├── members.json
├── metrics.json
├── profile.prof
└── trace.out
```
