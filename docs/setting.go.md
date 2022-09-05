# Setting `Go` on WSL

This post covers how to set [`Go`](https://go.dev/) on **Windows WSL2**.  

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

### [Go in `Visual Studio Code`](https://code.visualstudio.com/docs/languages/go)  
You can install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) from the VS Code Marketplace.  

> Watch ["Getting started with VS Code Go"](https://www.youtube.com/watch?v=1MXIGYrMk80) for an explanation of  
> how to build your first Go application using VS Code Go.  

<br/>

### [vscode-go Debugging Configuration](https://github.com/golang/vscode-go/blob/master/docs/debugging.md#configuration)
  ![dlv-dap](https://github.com/golang/vscode-go/raw/master/docs/images/vscode-go-debug-arch.png)  

* [Launch.json Attributes](https://github.com/golang/vscode-go/blob/master/docs/debugging.md#launchjson-attributes)
 
  Covers the list of attributes specific to Go debugging.  

  * *Launch*  
    This feature uses a launch request type configuration. Its program attribute needs to be either the go file or folder of the main package or test file. In this mode, the Go extension will start the debug session by building and launching the program. The launched program will be terminated when the debug session ends.  
  * *Attach*  
    You can use this configuration to attach to a running process or a running debug session.

* [Remote Debugging](https://github.com/golang/vscode-go/blob/master/docs/debugging.md#remote-debugging)  
  * [Connecting to Headless Delve with Target Specified at Server Start-Up](https://github.com/golang/vscode-go/blob/master/docs/debugging.md#connecting-to-headless-delve-with-target-specified-at-server-start-up)  
  * [Connecting to Delve DAP with Target Specified at Client Start-Up](https://github.com/golang/vscode-go/blob/master/docs/debugging.md#connecting-to-delve-dap-with-target-specified-at-client-start-up)  

<br/><br/><br/>

## References  
* [WSL2](https://docs.microsoft.com/en-us/windows/wsl/)  
* [Install Go](https://golang.org/doc/install)  
* [Go in `Visual Studio Code`](https://code.visualstudio.com/docs/languages/go)  
* [vscode-go Debugging](https://github.com/golang/vscode-go/blob/master/docs/debugging.md#connecting-to-headless-delve-with-target-specified-at-server-start-up)  