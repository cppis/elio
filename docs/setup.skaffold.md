# Setup `Skaffold`  
This post shows how to set [`Skaffold`](https://skaffold.dev/) up on **Windows 10**.  

<br/><br/>

## Prerequisites  
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

### [`kubectl`](https://kubernetes.io/docs/tasks/tools/)  
The Kubernetes command-line tool, `kubectl`, allows you to run commands  
against Kubernetes clusters. 

* [Install kubectl on Windows](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/)  

<br/>

### [`minikube`](https://minikube.sigs.k8s.io/docs/start/)   
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
          <img src="https://miro.medium.com/max/665/1*xwFelgX0H_c91tBknDu-_w.png" style="width: 320px; max-width: 100%; height: auto" title="cloudflare fixed the bug" />
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

<br/><br/><br/>


## References  
* [Install Hyper-V on Windows 10](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/quick-start/enable-hyper-v)  