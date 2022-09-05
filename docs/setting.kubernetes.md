# Setting `Kubernetes` on WSL

This post covers how to set `Kubernetes` + [`Skaffold`](https://skaffold.dev/) on **Windows WSL2**.  

<br/><br/>

## Installation  
### [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
The Kubernetes command-line tool, `kubectl`, allows you to run commands  
against Kubernetes clusters. 

* [Install and Set Up kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)  

<br/>

### [`Helm`](https://helm.sh/)  
`Helm` is a package manager for Kubernetes  
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

### [`Kind`](https://kind.sigs.k8s.io/)  
[`Kind`](https://kind.sigs.k8s.io/) is a tool for running local Kubernetes clusters using Docker container `nodes`.  
* [Kind Installation](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)  
  ```bash
  curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
  chmod +x ./kind
  sudo mv ./kind /usr/local/bin/kind
  ```

* [Setting `Kind Registry`](https://kind.sigs.k8s.io/docs/user/local-registry/#create-a-cluster-and-registry)  
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

    > If the Kind cluster does not have access to the local Docker registry,
    > Associate the Kind cluster's network with the local Docker registry's network by running the following command:  
    > ```bash
    > docker network connect "kind" "kind-registry"
    > ```

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

<br/><br/><br/>

## References  
* [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
* [`Helm`](https://helm.sh/)  
* [`Kind`](https://kind.sigs.k8s.io/)  
* [`Skaffold`](https://skaffold.dev/docs/install/)  
