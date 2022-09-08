# Troubleshooting  

<br/><br/><br/>

## Tips  

### 🧰 Run *docker cli* pod in Kubernetes:

```bash
kubectl run -it --rm --restart=Never dind --image=docker:dind -- sh
```

<br/>

### 🧰 Run *telnet* pod in Kubernetes:

```bash
kubectl run -it --rm --restart=Never busybox --image=gcr.io/google-containers/busybox -- sh
```

<br/>

### 🧰 Port forwarding a pod in Kubernetes:
 
```bash
kubectl port-forward $(kubectl get pods --selector=app=herald --output=jsonpath={.items..metadata.name}) 7003:7003
```

<br/>

### 🧰 Check kind cluster nodes:  

```bash
docker exec -it elio-control-plane crictl info
docker exec -it elio-control-plane crictl images
docker exec -it elio-worker crictl info
docker exec -it elio-worker crictl images
```

<br/><br/><br/>

## Issues  

## 🧰 Skaffold `ErrImagePull` error  
```
Failed to pull image "localhost:5001/skaffold-herald": 
  rpc error: code = Unknown desc = failed to pull and unpack image "localhost:5001/skaffold-herald:latest": 
  failed to resolve reference "localhost:5001/skaffold-herald:latest": failed to do request: 
  Head "http://localhost:5001/v2/skaffold-herald/manifests/latest": dial tcp [::1]:5001: connect: connection refused    
```

### References  
* [Pull-through Docker registry on Kind clusters](https://maelvls.dev/docker-proxy-registry-kind/)  
* [Skaffold: Tag](https://skaffold.dev/docs/pipeline-stages/taggers/)  
  the `sha256` tagger uses `latest`.  
* [KiND: How I Wasted a Day Loading Local Docker Images](https://iximiuz.com/en/posts/kubernetes-kind-load-docker-image/)  
<br/><br/><br/>
