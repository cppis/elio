# herald
`herald` means a messanger.  
It is a sample app to test pub/sub of backing `emqx` message broker  
using [`elio`](https://github.com/cppis/elio) library.  

If you run `skaffold` with `k8s-resources/profiles/local.yaml`,  
it is configured as follows in a Kubernetes:  

![docs/images/herald.skaffold.png](https://github.com/cppis/elio/blob/dev/docs/images/herald.skaffold.config.png?raw=true)  

<br/><br/><br/>

## Start  
### Create Kind  
```bash
kind create cluster --config app/herald/k8s-resources/kind-config.yaml --name elio
```

<br/>

### [Helm](https://helm.sh/)  
Helm is a package manager for Kubernetes  
Install from script:  
```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

Add `emqx` helm repo:  
```bash
helm repo add emqx https://repos.emqx.io/charts
```

<br/><br/><br/>

## Run herald  
### using `Skaffold`  
To use the `Skaffold`, you need thd following the [Setup `Skaffold`](#setup-skaffold).  
To run `herald` using `Skaffold`,  
run the following command in the Project root directory:  
```shell
$ skaffold -f app/herald/k8s-resources/skaffold.yaml dev
```

> To change detection triggered to manual mode, use option `--trigger=manual`.  

Or, to run `herald` in debugging mode using `Skaffold`, run the following command:  
```shell
$ skaffold -f app/herald/k8s-resources/skaffold.yaml debug
```

<br/>

### using `Docker`  
```bash
docker build --no-cache -t localhost:5001/skaffold-herald:latest -f app/herald/k8s-resources/Dockerfile .
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

## Debugging  

* run _telnet_ using _busybox_:

```bash
kubectl run -it --rm --restart=Never busybox --image=gcr.io/google-containers/busybox -- sh
```

* port forwarding pod by selector:
 
```bash
kubectl port-forward $(kubectl get pods --selector=app=herald --output=jsonpath={.items..metadata.name}) 7000:7000
```

<br/><br/><br/>

## Stop  
### Delete Kind  
```bash
kind delete clusters elio
```

<br/><br/><br/>

## Reference  
* [Building K8S cluster of EMQ X starting from scratch](https://www.emqx.com/en/blog/emqx-mqtt-broker-k8s-cluster)  


<br/><br/><br/>

## TO-DO  
* Run `emqx` as stateful cluster in kubernetes  
