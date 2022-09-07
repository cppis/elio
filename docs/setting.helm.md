# Setting `Helm` on Windows

This post shows how to set [`Helm`](https://skaffold.dev/) up on **Windows 10**.  

<br/><br/>

## [Installing `Helm`](https://helm.sh/docs/intro/install/)  

#### From the Binary Releases  
Every release of Helm provides binary releases for a variety of OSes.  
These binary versions can be manually downloaded and installed.

1. [Download your desired version](https://github.com/helm/helm/releases)  
2. Unpack it (`tar -zxvf helm-v3.0.0-linux-amd64.tar.gz`)  
3. Find the helm binary in the unpacked directory,  
   and move it to its desired destination  
   (`mv linux-amd64/helm /usr/local/bin/helm`)  

<br/>

### From Script  
Helm now has an installer script that will automatically grab  
the latest version of Helm and [install it locally](https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3).  

You can fetch that script, and then execute it locally. It's well documented so that you can read through it and understand what it is doing before you run it.   

```shell
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

Yes, you can `curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash` if you want to live on the edge.  

<br/><br/><br/>

## References  
* [Helm](https://helm.sh/)  