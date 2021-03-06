# Setting `Skaffold` on Windows

This post shows how to set [`Skaffold`](https://skaffold.dev/) up on **Windows 10**.  

<br/><br/>

## Prerequisites  
### Windows 10+
To update to WSL 2, you must be running Windows 10 or higher.

* For x64 systems: **Version 1903** or higher, with **Build 18362** or higher.
* For ARM64 systems: **Version 2004** or higher, with **Build 19041** or higher.
* Builds lower than **18362** do not support WSL 2. Use the [Windows Update Assistant](https://www.microsoft.com/ko-kr/software-download/windows10ISO) to update your version of Windows.

<br/>

### [Docker Desktop](https://www.docker.com/products/docker-desktop)  
Docker Desktop includes a standalone Kubernetes server and client, as well as Docker CLI integration that runs on your machine. The Kubernetes server runs locally within your Docker instance, is not configurable, and is a single-node cluster.  

* [Enable Kubernetes](https://docs.docker.com/desktop/kubernetes/#enable-kubernetes)  
  To enable Kubernetes support and install a standalone instance of Kubernetes running as a Docker container, go to Preferences > Kubernetes and then click Enable Kubernetes.  

  <figure>
    <div style="text-align:center">
      <img src="https://docs.docker.com/desktop/images/kube-enable.png" style="width: 480px; max-width: 100%; height: auto" title="cloudflare fixed the bug" />
    </div>
  </figure>

<br/>

### [`Go` 1.16+](https://golang.org/doc/install)  
Go is an open source programming language that makes it easy  
to build simple, reliable, and efficient software.  
[Download](https://golang.org/doc/install#download) and [Install](https://golang.org/doc/install#install) [`Go`](https://golang.org/) v1.6 or higher for Windows   

  <figure>
  <div style="text-align:center">
    <a href="https://drive.google.com/uc?export=view&id=13pf8blQenwo2SBnkebnMkg9vai5f8REa">
    <img src="https://drive.google.com/uc?export=view&id=13pf8blQenwo2SBnkebnMkg9vai5f8REa" style="width: 100px; max-width: 100%; height: auto" title="hypver-v-virtual-switches" />
    </a>
  </div>
  </figure>

<br/>

### [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
The Kubernetes command-line tool, `kubectl`, allows you to run commands  
against Kubernetes clusters. 

* [Install kubectl on Windows](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/)  

<br/>

### [`minikube`(Optional)](https://minikube.sigs.k8s.io/docs/start/)   
`minikube` quickly sets up a local Kubernetes cluster on macOS, Linux, and Windows.  
`minikube` focus on helping application developers and new Kubernetes users.  

* Download and run the stand-alone minikube [Windows installer](https://storage.googleapis.com/minikube/releases/latest/minikube-installer.exe).  
* [hyperv driver](https://minikube.sigs.k8s.io/docs/drivers/hyperv/)  
  [Hyper-V](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/) is a native hypervisor built in to modern versions of Microsoft Windows.  
  * [Enabling Hyper-V](https://minikube.sigs.k8s.io/docs/drivers/hyperv/#enabling-hyper-v)  
  <figure>
    <div style="text-align:center">
      <img src="https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/quick-start/media/enable_role_upd.png" style="width: 320px; max-width: 100%; height: auto" title="cloudflare fixed the bug" />
    </div>
  </figure>


    ```shell
    Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
    ```
    If Hyper-V was not previously active, you will need to **reboot**.  

  * Add Virtual Network Switch  
    * Type `hyper-v` in search bar to open the **Hyper-V Manager** and run it.  
    * Select `Virtual Switch Manager...` on the right panel.  
    * Select `New virtual network switch` on the left panel, select `External`  
      and press `Create Virtual Switch` button to create a virtual switch for minikube.  

      <figure>
      <div style="text-align:center">
        <a href="https://drive.google.com/uc?export=view&id=1NKLzsTq1L3s8bL0-worZnS6rGfL4iZey">
        <img src="https://drive.google.com/uc?export=view&id=1NKLzsTq1L3s8bL0-worZnS6rGfL4iZey" style="width: 500px; max-width: 100%; height: auto" title="hypver-v-virtual-switches" />
        </a>
      </div>
      </figure>

    * Name the switch **Primary Virtual Switch** and click the **OK** button.  
    * To make *hyperv* the default driver:
        ```shell
        minikube config set driver hyperv
        ```

    <br/>

    > Once you have the switch created we are now ready to start minikube.  
    > If you don't want to use *hyperv* for default,  
    > Run the following command to start the minikube VM with our applied changes.  
    > 
    > ```shell
    > minikube start --vm-driver hyperv --hyperv-virtual-switch "Primary Virtual Switch"  
    > ```

<br/>

### [`Skaffold`](https://skaffold.dev/docs/install/)  
`Skaffold` is a command line tool that facilitates continuous development for Kubernetes-native applications. Skaffold handles the workflow for building, pushing, and deploying your application, and provides building blocks for creating CI/CD pipelines. This enables you to focus on iterating on your application locally while Skaffold continuously deploys to your local or remote Kubernetes cluster.  

* [Download Skaffold](https://storage.googleapis.com/skaffold/releases/latest/skaffold-windows-amd64.exe)  
  The latest stable release binary can be found here:  
  Simply download it and place it in your `PATH` as `skaffold.exe`.

  > You can permanently add a path to system `PATH`:   
  > ```shell
  > $ setx /M path "%path%;{Skaffold Path}"
  > ```
  > Or add to user `PATH`:   
  > ```shell
  > $ setx path "%path%;{Skaffold Path}"
  > ```


#### [`--generate-manifests` Flag](https://skaffold.dev/docs/pipeline-stages/init/#--generate-manifests-flag)  
`beta`  
`skaffold init` allows for use of a `--generate-manifests` flag, which will try 
to generate basic kubernetes manifests for a user???s project to help get things up and running. If bringing a project to skaffold that has no kubernetes manifests yet, it may be helpful to run skaffold init with this flag.  

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
(Assumes `minikube` default driver is hyperv)
```shell
minikube start
cd app/echo
```

`skaffold dev` enables continuous local development on an application.  
```shell
skaffold dev -p dev
```

Or, `skaffold debug` acts like `skaffold dev`, but it configures containers in pods  
for debugging as required for each container???s runtime technology.  
```shell
skaffold debug -p debug
```

<br/>

### Stop  
The dev loop will run until the user cancels the `Skaffold` process with `Ctrl+C`.  
Upon receiving this signal, `Skaffold` will clean up all deployed artifacts on the active cluster.  
This can be optionally disabled by using the `--no-prune` flag.  

<br/>

To stop minikube, run the following command:  
```shell
minikube stop
```

<br/>

### Delete  
To delete minikube, run the following command:  
```shell
minikube delete
```

<br/><br/><br/>

## References  
* [Install Hyper-V on Windows 10](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/quick-start/enable-hyper-v)  