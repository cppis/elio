# Troubleshooting  

Shows some tips to solve the problem.  

<br/><br/>

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

The issue was caused by failing to fetch images while running the `Scaffold`, the problem was with the 'Skaffold' *tagPolicy* setting.  
When you use the `latest` tag, you should use the `sha256` *tagPolicy*.  

<br/>

### References  
* [Pull-through Docker registry on Kind clusters](https://maelvls.dev/docker-proxy-registry-kind/)  
* [Skaffold: Tag](https://skaffold.dev/docs/pipeline-stages/taggers/)  
  the `sha256` tagger uses `latest`.  
* [KiND: How I Wasted a Day Loading Local Docker Images](https://iximiuz.com/en/posts/kubernetes-kind-load-docker-image/)  

<br/><br/><br/>

## Create [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) cluster error  
```bash
$ ./app/assets.k8s/kind.with.registry.sh 
Creating cluster "elio" ...
 ✓ Ensuring node image (kindest/node:v1.24.0) 🖼 
 ✓ Preparing nodes 📦 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✗ Joining worker nodes 🚜 
ERROR: failed to create cluster: failed to join node with kubeadm: command "docker exec --privileged elio-worker kubeadm join --config /kind/kubeadm.conf --skip-phases=preflight --v=6" failed with error: exit status 1
Command Output: I0916 01:16:36.105116     303 join.go:413] [preflight] found NodeName empty; using OS hostname as NodeName
I0916 01:16:36.105146     303 joinconfiguration.go:76] loading configuration from "/kind/kubeadm.conf"
I0916 01:16:36.106267     303 controlplaneprepare.go:220] [download-certs] Skipping certs download
I0916 01:16:36.106273     303 join.go:530] [preflight] Discovering cluster-info
```

*Joining worker nodes* error may occur when creating multiple Kind clusters.  
This may be caused by running out of [inotify](https://linux.die.net/man/7/inotify) resources. Resource limits are defined by `fs.inotify.max_user_watches` and `fs.inotify.max_user_instances` system variables.  

This can be solved by fixing the `ulimit` settings:

```bash
echo fs.inotify.max_user_watches=655360 | sudo tee -a /etc/sysctl.conf
echo fs.inotify.max_user_instances=1280 | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```
