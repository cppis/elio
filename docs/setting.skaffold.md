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

### [`Go` 1.19+](https://golang.org/doc/install)  
Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.  
[Download](https://golang.org/doc/install#download) and [Install](https://golang.org/doc/install#install) [`Go`](https://golang.org/) v1.19 or higher  

<br/>

### [Docker Desktop](https://www.docker.com/products/docker-desktop)  
*Docker Desktop* is an application for MacOS and Windows machines for the building and sharing of containerized applications and microservices.  

<br/>

### [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
The Kubernetes command-line tool, `kubectl`, allows you to run commands  
against Kubernetes clusters. 

* [Install and Set Up kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)  

<br/>

### [Helm](https://helm.sh/)  
Helm is a package manager for Kubernetes  
Install from script:  
```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

After installation, Add helm repo.  
For example, you can add a `emqx` repo with the following command:

```bash
helm repo add emqx https://repos.emqx.io/charts
```

<br/>

### [Kind](https://kind.sigs.k8s.io/)  
[Kind](https://kind.sigs.k8s.io/) is a tool for running local Kubernetes clusters using Docker container `nodes`.  
* [Kind Installation](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)  
  ```bash
  curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
  chmod +x ./kind
  sudo mv ./kind /usr/local/bin/kind
  ```

* [Create A Cluster And Registry](https://kind.sigs.k8s.io/docs/user/local-registry/#create-a-cluster-and-registry)
  If you don't have a Kubernetes local image registry, create one.  
  First, check if you already have a local docker registry with the following command:  

  ```bash
  # Make sure kind-registry is running.  
  docker inspect -f '{{.State.Running}}' kind-registry
  ```

  If the result of the above command is not `True` , create a local docker registry with the following command:  

  ```bash
  docker run -d --restart=always -p 127.0.0.1:5001:5000 --name kind-registry registry:2
  ```

  After the installation, you can check the docker registry with the following command:  

  ```bash
  curl http://127.0.0.1:5001/v2/_catalog
    {"repositories":[]}
  ```

  > For more information, refer to [Kind - Local Registry](https://kind.sigs.k8s.io/docs/user/local-registry/).  

<br/>

### [`Skaffold`](https://skaffold.dev/docs/install/)  
`Skaffold` is a command line tool that facilitates continuous development for Kubernetes-native applications. Skaffold handles the workflow for building, pushing, and deploying your application, and provides building blocks for creating CI/CD pipelines. This enables you to focus on iterating on your application locally while Skaffold continuously deploys to your local or remote Kubernetes cluster.  

* [Installing Skaffold](https://skaffold.dev/docs/install/)  
  ```
  curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
  sudo install skaffold /usr/local/bin/
  ```

The following figure shows `Skaffolder` workflow:

<figure>
  <div style="text-align:center">
    <img src="https://skaffold.dev/images/architecture.png" style="width: 640px; max-width: 100%; height: auto" title="cloudflare fixed the bug" />
  </div>
</figure>

<br/>

### [Go in `Visual Studio Code`](https://code.visualstudio.com/docs/languages/go)  
You can install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) from the VS Code Marketplace.  

> Watch ["Getting started with VS Code Go"](https://www.youtube.com/watch?v=1MXIGYrMk80) for an explanation of  
> how to build your first Go application using VS Code Go.  

<br/><br/><br/>

## References  
* [WSL2](https://docs.microsoft.com/en-us/windows/wsl/)  
* [Install Go](https://golang.org/doc/install)  
* [Docker Desktop](https://www.docker.com/products/docker-desktop)  
* [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
* [Helm](https://helm.sh/)  
* [Kind](https://kind.sigs.k8s.io/)  
* [`Skaffold`](https://skaffold.dev/docs/install/)  
* [Go in `Visual Studio Code`](https://code.visualstudio.com/docs/languages/go)  
