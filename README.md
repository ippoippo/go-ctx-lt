# go-ctx-lt

Source code that will be used in a blog post.

## Pre-requisites

- Go [1.23.1 or later](https://go.dev/doc/devel/release#go1.23.minor) - [Download](https://go.dev/dl/)
- Optional: [Task](https://taskfile.dev) : `make`, but via yaml specifications
  - `brew install go-task`

## Tasks

- `task timeout-demo` : Simple demonstration of a timeout on an API call.
- `task timeout-with-cause-demo` : Adds in custom error for when a timeout occurs.
- `task deadline-demo` : Based on `timeout-demo`, but shows we can use Deadline instead.
- `task statuschecker-demo` : Demonstrates use of context cancel.
- `task simplecontextvalue-demo` : Demonstrates simplistic context values.
- `task contextvalue-demo` : Demonstrates using context values, with package keys, to set trace-ids in a webserver.
- `task contextvalue-demo-curl` : To be used in connjuction with `contextvalue-demo`.
- `task parentchildcancellation-demo` : Demonstrates using that cancellation only propogates down
