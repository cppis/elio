# Herald

`Herald` means a messanger.  
It is a simple app to test pub/sub of backing `emqx` MQTT broker.  
`Herald` is written using [`elio`](https://github.com/cppis/elio) library.  

![herald.concept](https://github.com/cppis/elio/blob/dev/docs/images/herald.concept.png?raw=true)  

* Docker Hub: [cppis/herald](https://hub.docker.com/repository/docker/cppis/herald)  

<br/><br/><br/>

## Prerequisites  

### Setting up [`elio`](https://github.com/cppis/elio)  

Before start, set up `elio` project.  
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

### [Setting `Kubernetes` on WSL](https://github.com/cppis/elio/blob/dev/docs/setting.kubernetes.md)  

Covers `Kubernetes`+`Skaffold` settings on WSL for continuous developing a Kubernetes-native app.  


<br/>

### [Setting Service Ports](https://github.com/cppis/elio/blob/dev/docs/setting.service.ports.md)  

Covers the port settings for the Kubernetes *Service*.  

<br/><br/><br/>

## Running app on Host  

You can easily run a `Herald` container on the host (without MQTT Broker):  

![herald.container](https://github.com/cppis/elio/blob/dev/docs/images/herald.container.png?raw=true)  

<br/>

### Using `go run`  
To run `Herald` service, run the following command:  
```shell
HERALD_IN_URL="0.0.0.0:7002" go run ./app/herald
```

You can change the listening url of service `Herald` by changing  
environment variable `HERALD_IN_URL`.

<br/>

### Using `Docker`  
To run `Herald` container, run the following command:  
```bash
docker run -d -e HERALD_IN_URL="0.0.0.0:7002" -p 7002:7002 --name herald cppis/herald:latest
```

To kill `Herald` container, run the following command:  
```bash
docker rm -f herald 
```

<br/><br/><br/>

## Running app on Kubernetes  

You can easily run a `Herald` + `emqx`(MQTT Broker) chart on the kubernetes cluster.   

![herald.chart](https://github.com/cppis/elio/blob/dev/docs/images/herald.chart.png?raw=true)  

<br/>

If the kubernetes cluster does not exist, Follow the next step `Create a Kind cluster`.  
If you have, Follow the next step 'Using `Helm`' or 'Using `Skaffold`'.  

<br/>

### Create a `Kind` cluster  

To create a `kind`, run the following command:  
```bash
kind create cluster --config app/herald/assets.k8s/kind.cluster.yaml --name elio
```

> To check if the kind cluster is up and running, run the following command:  
> ```bash
> kubectl config current-context
>   kind-elio
> ```

<br/>

### Using `Helm`  

To deploy using helm, run the following command:  
```bash
helm upgrade --install herald app/herald/assets.k8s/helm
```

helm uninstall herald

<br/>

### Using `Skaffold`  

To run `herald` using `Skaffold`,  
run the following command in the Project root directory:  
```shell
skaffold -f app/herald/assets.k8s/skaffold.yaml dev
```

> To debugging `Skaffold`, use option `-vdebug`.  

> To change detection triggered to manual mode, use option `--trigger=manual`.  

Or, to run `Herald` in debugging mode using `Skaffold`, run the following command:  
```shell
skaffold -f app/herald/assets.k8s/skaffold.yaml debug
```

<br/><br/><br/>

## Testing app  
You can test echo easily by using telnet.  

app protocol is custom `t2p` like http.  
procotol header is separated by newline(`\n` or `\r\n`).  
And packet delimiter is double newline(`\n\n` or `\r\n\r\n`).

<br/>

### connect
connect to echo using `telnet`  
  ```bash
  telnet localhost 7002
  ```

<br/>

### echo
echo message    
  ```
  echo<newline>
  {message}<newline><newline>
  ```

<br/>

### sub
subcribe to topic    
  ```
  sub<newline>
  {topic}<newline><newline>
  ```

<br/>

### unsub
unsubcribe from topic  
  ```
  unsub<newline>
  {topic}<newline><newline>
  ```

<br/>

### pub: publish message to topic  
  ```
  pub<newline>
  {topic}<newline>
  {message}<newline><newline>
  ```

<br/><br/><br/>

## Debugging Tips  

### Running _telnet_ in Kubernetes:

```bash
kubectl run -it --rm --restart=Never busybox --image=gcr.io/google-containers/busybox -- sh
```

<br/>

### Port forwarding a pod in Kubernetes:
 
```bash
kubectl port-forward $(kubectl get pods --selector=app=herald --output=jsonpath={.items..metadata.name}) 7002:7002
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
