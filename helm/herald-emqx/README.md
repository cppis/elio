# herald with helm

build and tag image from Project Root:  
```shell
$ cd {Project Root}
$ docker build -t herald:v0.0.1 -f /app/herald/k8s-resources/Dockerfile .
```

Test docker image:  
```shell
$ docker run -it -p 7000:7000 -e HERALD_IN_URL=0.0.0.0:7000 herald:v0.0.1
  ...
```

<br/><br/>

## Create a *chart*  
To create chart:  
```shell
$ helm create herald
```

<br/><br/>

## Install a *release*  
To install chart as release `herald-emqx-release`:  
```shell
$ helm install 
    -f hello-world/values.yaml
    -n herald
    herald-emqx-release helm/herald-emqx
```

helm install herald-emqx-release helm/herald-emqx

> To pass environment variable add `--set` option like `--set HERALD_IN_URL="0.0.0.0:800"`

To find exposed service port:  
```shell
$ kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services herald-emqx-release
$ kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}"
```

To list release:  
```shell
helm list
```

<br/><br/>

## Uninstall a *release*  
To uninstall a release `herald-emqx-release`:  
```shell
$ helm uninstall herald-emqx-release
```

```shell
helm delete herald-emqx-release --purge
```

delete namespace in the kubernetes cluster:  
```shell
$ kubectl delete ns herald
```
