## How to run

```shell
go mod download
export DO_TOKEN=<some_token_here>
export DO_CLUSTER_ID=<k8s_cluster_here>
export K8S_NAMESPACE=<namespace_here>
go run main.go -v=8
```
