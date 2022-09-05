# Setting `Docker` on WSL

This post covers how to set [[`Docker`](https://skaffold.dev/)](https://www.docker.com/products/docker-desktop) on **Windows WSL2**.  

<br/><br/>

## Installation  
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

### [`Docker Desktop`](https://www.docker.com/products/docker-desktop)  
`Docker Desktop` is an application for MacOS and Windows machines for the building and sharing of containerized applications and microservices.  

<br/><br/><br/>

## References  
* [WSL2](https://docs.microsoft.com/en-us/windows/wsl/)  
* [Install Go](https://golang.org/doc/install)  
* [Docker Desktop](https://www.docker.com/products/docker-desktop)  
