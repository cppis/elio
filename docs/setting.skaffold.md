# Setting `Skaffold` on WSL

This post shows how to set [`Skaffold`](https://skaffold.dev/) up on **Windows WSL2**.  

<br/><br/>

## Prerequisites  
### [WSL2](https://docs.microsoft.com/en-us/windows/wsl/)  
The Windows Subsystem for Linux(WSL) lets developers run a GNU/Linux environment -- including most command-line tools, utilities, and applications -- directly on Windows, unmodified, without the overhead of a traditional virtual machine or dualboot setup.  

* [Install Linux on Windows with WSL](https://docs.microsoft.com/en-us/windows/wsl/install)
* [Advanced settings configuration in WSL](https://docs.microsoft.com/en-us/windows/wsl/wsl-config)

> Note: To update to WSL 2, you must be running Windows 10 or higher.  
> * For x64 systems: **Version 1903** or higher, with **Build 18362** or higher.
> * For ARM64 systems: **Version 2004** or higher, with **Build 19041** or higher.
> * Builds lower than **18362** do not support WSL 2. Use the [Windows Update Assistant](https://www.microsoft.com/ko-kr/software-download/windows10ISO) to update your version of Windows.

<br/>

### [Docker Desktop](https://www.docker.com/products/docker-desktop)  
*Docker Desktop* is an application for MacOS and Windows machines for the building and sharing of containerized applications and microservices.  

<br/>

### [`Go` 1.19+](https://golang.org/doc/install)  
Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.  
[Download](https://golang.org/doc/install#download) and [Install](https://golang.org/doc/install#install) [`Go`](https://golang.org/) v1.19 or higher  

<br/>

### [Kind](https://kind.sigs.k8s.io/)  
[Kind](https://kind.sigs.k8s.io/) is a tool for running local Kubernetes clusters using Docker container `nodes`.  
* [Kind Installation](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)  
  ```bash
  curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
  chmod +x ./kind
  sudo mv ./kind /usr/local/bin/kind
  ```

<br/>

### [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
The Kubernetes command-line tool, `kubectl`, allows you to run commands  
against Kubernetes clusters. 

* [Install and Set Up kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)  

<br/>

### [`Skaffold`](https://skaffold.dev/docs/install/)  
`Skaffold` is a command line tool that facilitates continuous development for Kubernetes-native applications. Skaffold handles the workflow for building, pushing, and deploying your application, and provides building blocks for creating CI/CD pipelines. This enables you to focus on iterating on your application locally while Skaffold continuously deploys to your local or remote Kubernetes cluster.  

* [Installing Skaffold](https://skaffold.dev/docs/install/)  
  ```
  curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
  sudo install skaffold /usr/local/bin/
  ```

<br/>

### [Go in `Visual Studio Code`](https://code.visualstudio.com/docs/languages/go)  
You can install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) from the VS Code Marketplace.  

> Watch ["Getting started with VS Code Go"](https://www.youtube.com/watch?v=1MXIGYrMk80) for an explanation of  
> how to build your first Go application using VS Code Go.  

<br/><br/><br/>

## Run Echo  
The following figure shows `Skaffolder` workflow.  

<figure>
  <div style="text-align:center">
    <img src="https://skaffold.dev/images/architecture.png" style="width: 640px; max-width: 100%; height: auto" title="cloudflare fixed the bug" />
  </div>
</figure>

### Start  
After all settings are done, run local kubernetes by the following command:  
```shell
kind create cluster --config app/echo/k8s-resources/kind-config.yaml --name elio
kubectl config current-context
  kind-elio
cd $ELIO_ROOT
```

`skaffold dev` enables continuous local development on an application.  
```shell
skaffold -f ./app/echo/k8s-resources/skaffold.yaml dev -p dev
```

Or, `skaffold debug` acts like `skaffold dev`, but it configures containers in pods  
for debugging as required for each container’s runtime technology.  
```shell
skaffold -f ./app/echo/k8s-resources/skaffold.yaml debug -p debug
```

<br/>

### Test  
You can test echo easily by using telnet.  
And, you can end server by send `Ctrl+c`.  

```
telnet localhost 7000
```

<br/>

### Stop  
#### [Kubernetes resource cleanup](https://skaffold.dev/docs/pipeline-stages/cleanup/#kubernetes-resource-cleanup)  
After running `skaffold run` or `skaffold deploy` and deploying your app to a cluster, running `skaffold delete` will remove all the resources you deployed. Cleanup is enabled by default, it can be turned off by `--cleanup=false`  

#### [Ctrl + C](https://skaffold.dev/docs/pipeline-stages/cleanup/#ctrl--c)  
When running `skaffold dev` or `skaffold debug`, pressing `Ctrl+C` (SIGINT signal) will kick off the cleanup process which will mimic the behavior of `skaffold delete`. If for some reason the Skaffold process was unable to catch the SIGINT signal, `skaffold delete` can always be run later to clean up the deployed Kubernetes resources.

To enable image pruning, you can run Skaffold with both `--no-prune=false` and `--cache-artifacts=false`:

```
skaffold dev --no-prune=false --cache-artifacts=false
```

<br/>

### Delete  
To delete kind, run the following command:  
```shell
kind delete clusters elio
```

<br/><br/><br/>

## References  
* [Install Hyper-V on Windows 10](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/quick-start/enable-hyper-v)  