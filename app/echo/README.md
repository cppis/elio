# Echo

`Echo` is a simple echo server using elio library.  

<br/><br/><br/>

## Prerequisites  

### Setting up [`elio`](https://github.com/cppis/elio)  

Before start, set up `elio` project:  
```
git clone https://github.com/cppis/elio && cd elio
go mod vendor
export ELIO_ROOT=$(pwd)
```

> Now, **$ELIO_ROOT** is the project root path.  

<br/>

### [Setting `Docker` on WSL](https://github.com/cppis/elio/blob/dev/docs/setting.docker.md)  

Covers `Docker` settings on WSL.  

<br/>

### [Setting `Kubernetes`](https://github.com/cppis/elio/blob/dev/docs/setting.kubernetes.md)  

Covers `Kubernetes`+`Skaffold` settings on WSL for continuous developing a Kubernetes-native app.  

<br/><br/><br/>

## Running app on Host  
### Using `go run`  
To run `echo` service, run the following command:  
```shell
ECHO_IN_URL="0.0.0.0:7001" go run ./app/echo
```

You can change the listening url of service `echo` by changing  
environment variable `ECHO_IN_URL`.

<br/><br/><br/>

## Running app on Kubernetes  

If the kubernetes cluster does not exist, Follow the next step `Create a Kind cluster`.  
If you have, Follow the next step Using `Skaffold`.  

### Create a Kind cluster  

To create a kind, run the following command:  
```bash
kind create cluster --config app/echo/assets.k8s/kind.cluster.yaml --name elio
```

> To check if the kind cluster is up and running, run the following command:  
> ```bash
> kubectl config current-context
>   kind-elio
> ```

<br/>

### Using `Skaffold`  

To run `echo` using `Skaffold`,  
run the following command:  
```shell
skaffold -f app/echo/assets.k8s/skaffold.yaml dev -p dev
```

> To debugging `Skaffold`, use option `-vdebug`.  

Or, `skaffold debug` acts like `skaffold dev`, but it configures containers in pods  
for debugging as required for each container’s runtime technology.  
```shell
skaffold -f app/echo/assets.k8s/skaffold.yaml debug -p debug
```

<br/><br/><br/>

## Testing app  
You can test echo easily by using telnet.  
And, you can end server by send `q` character.  

```
telnet localhost 7001
  ...
  q<enter>
```

<br/><br/><br/>

## Ending app  
### [Kubernetes resource cleanup](https://skaffold.dev/docs/pipeline-stages/cleanup/#kubernetes-resource-cleanup)  
After running `skaffold run` or `skaffold deploy` and deploying your app to a cluster, running `skaffold delete` will remove all the resources you deployed. Cleanup is enabled by default, it can be turned off by `--cleanup=false`  

<br/>

### [Ctrl + C](https://skaffold.dev/docs/pipeline-stages/cleanup/#ctrl--c)  
When running `skaffold dev` or `skaffold debug`, pressing `Ctrl+C` (SIGINT signal) will kick off the cleanup process which will mimic the behavior of `skaffold delete`. If for some reason the Skaffold process was unable to catch the SIGINT signal, `skaffold delete` can always be run later to clean up the deployed Kubernetes resources.

To enable image pruning, you can run Skaffold with both `--no-prune=false` and `--cache-artifacts=false`:

```
skaffold dev --no-prune=false --cache-artifacts=false
```

<br/>

### Delete a Kind cluster  
To delete a kind, run the following command:  
```bash
kind delete cluster --name elio
```
