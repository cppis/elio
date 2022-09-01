# Herald
`herald` means a messanger.  
It is a sample app to test pub/sub of backing `emqx` message broker  
using [`elio`](https://github.com/cppis/elio) library.  

If you run `skaffold` with `assets.k8s/skaffold.yaml`,  
it is configured as follows in a Kubernetes:  

![docs/images/herald.helm](https://github.com/cppis/elio/blob/dev/docs/images/herald.helm.png?raw=true)  

<br/><br/><br/>

## Prerequisites  
### [Setting `Skaffold` on WSL](docs/setting.skaffold.md)  

Before start, set the working path to the **$ELIO_ROOT**.  
```
cd $ELIO_ROOT
```
> **$ELIO_ROOT** is the project root path.  

<br/><br/><br/>

## Running on Host  
### Using `go run`  
To run `herald` service, run the following command:  
```shell
HERALD_IN_URL="0.0.0.0:7002" go run app/herald
```

You can change the listening url of service `herald` by changing  
environment variable `HERALD_IN_URL`.

<br/><br/><br/>

## Running on Kubernetes  
### Create a Kind cluster  

To create a kind, run the following command:  
```bash
kind create cluster --config app/herald/assets.k8s/kind.cluster.yaml --name elio
```

> To check if the kind cluster is up and running, run the following command:  
> ```bash
> kubectl config current-context
>   kind-elio
> ```

### Using the `Helm`  

To deploy using helm, run the following command:  
```bash
helm upgrade --install herald app/herald/assets.k8s/helm
```

<br/>

### Using the `Kind`+`Skaffold`  

To run `herald` using `Skaffold`,  
run the following command in the Project root directory:  
```shell
skaffold -f app/herald/assets.k8s/skaffold.yaml dev -vdebug
```

> To change detection triggered to manual mode, use option `--trigger=manual`.  

Or, to run `herald` in debugging mode using `Skaffold`, run the following command:  
```shell
skaffold -f app/herald/assets.k8s/skaffold.yaml debug
```

<br/>

### using `Helm`  
```bash
cd $ELIO_ROOT
helm install herald .
helm uninstall herald
```

<br/>

### using `Docker`  
```bash
docker build --no-cache -t localhost:5001/skaffold-herald:latest -f app/herald/assets.k8s/Dockerfile .
docker push localhost:5001/skaffold-herald:latest
```

<br/><br/><br/>

## Test  
You can test echo easily by using telnet.  

app protocol is custom `t2p` like http.  
procotol header is separated by newline(`\n` or `\r\n`).  
And packet delimiter is double newline(`\n\n` or `\r\n\r\n`).

### echo: echo message    
  ```
  echo<newline>
  {message}<newline><newline>
  ```
### sub: subcribe to topic    
  ```
  sub<newline>
  {topic}<newline><newline>
  ```
### unsub: unsubcribe from topic  
  ```
  unsub<newline>
  {topic}<newline><newline>
  ```
### pub: publish message to topic  
  ```
  pub<newline>
  {topic}<newline>
  {message}<newline><newline>
  ```

<br/><br/><br/>

## Debugging Tips  

* Running _telnet_ in Kubernetes:

```bash
kubectl run -it --rm --restart=Never busybox --image=gcr.io/google-containers/busybox -- sh
```

* Port forwarding a pod in Kubernetes:
 
```bash
kubectl port-forward $(kubectl get pods --selector=app=herald --output=jsonpath={.items..metadata.name}) 7002:7002
```

<br/><br/><br/>

<br/><br/><br/>

## Ending  
### [Kubernetes resource cleanup](https://skaffold.dev/docs/pipeline-stages/cleanup/#kubernetes-resource-cleanup)  
After running `skaffold run` or `skaffold deploy` and deploying your app to a cluster, running `skaffold delete` will remove all the resources you deployed. Cleanup is enabled by default, it can be turned off by `--cleanup=false`  

### [Ctrl + C](https://skaffold.dev/docs/pipeline-stages/cleanup/#ctrl--c)  
When running `skaffold dev` or `skaffold debug`, pressing `Ctrl+C` (SIGINT signal) will kick off the cleanup process which will mimic the behavior of `skaffold delete`. If for some reason the Skaffold process was unable to catch the SIGINT signal, `skaffold delete` can always be run later to clean up the deployed Kubernetes resources.

To enable image pruning, you can run Skaffold with both `--no-prune=false` and `--cache-artifacts=false`:

```
skaffold dev --no-prune=false --cache-artifacts=false
```

### Delete a Kind cluster  
To delete a kind, run the following command:  
```bash
kind delete cluster --name elio
```

<br/><br/><br/>

## Reference  
* [Building K8S cluster of EMQ X starting from scratch](https://www.emqx.com/en/blog/emqx-mqtt-broker-k8s-cluster)  


<br/><br/><br/>

## TO-DO  
* Run `emqx` as stateful cluster in kubernetes  
