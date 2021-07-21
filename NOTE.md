# Note  

## Prerequisite  
* [Docker](https://www.docker.com/)  
* [minikube](https://minikube.sigs.k8s.io/docs/)  
    + hyper-v
* [kubectl](https://kubernetes.io/docs/tasks/tools/)   

<br/><br/><br/>

## Configuration  
### Package Management  

```shell
$ skaffold build -vdebug  
```

<br/>

### Skaffold  

<br/><br/><br/>

## Issues  
* Dockerfile 내의 `go mod download` 시 timeout 이 발생함.  
    Usually the very first thing you do once you’ve downloaded a project written in Go is to install the modules necessary to compile it.

    But before we can run go mod download inside our image, we need to get our go.mod and go.sum files copied into it. We use the COPY command to do this.  
  * Multi-stage builds

<br/><br/><br/>

## References  
* [skaffold](https://skaffold.dev/)  

<br/><br/><br/>

## Troubleshooting  
* [Build your Go image](https://docs.docker.com/language/golang/build-images/)  
