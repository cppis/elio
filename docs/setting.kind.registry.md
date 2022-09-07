# Setting `Kind Registry`  

<br/>

## [Create A Cluster And Registry](https://kind.sigs.k8s.io/docs/user/local-registry/#create-a-cluster-and-registry)  

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
```

The result is as follows:  
```
{"repositories":[]}
```

> For more information, refer to [Kind - Local Registry](https://kind.sigs.k8s.io/docs/user/local-registry/).  

> If the Kind cluster does not have access to the local Docker registry,
> Associate the Kind cluster's network with the local Docker registry's network by running the following command:  
> ```bash
> docker network connect "kind" "kind-registry"
> ```
