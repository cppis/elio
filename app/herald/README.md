# herald
`herald` means messanger.  
It is a sample app to test pub/sub of backing `emqx` message broker.  

`k8s-resources/profiles/local.yaml`:  

![docs/images/herald.skaffold.png](https://github.com/cppis/elio/blob/dev/docs/images/herald.skaffold.config.png?raw=true)  

<br/><br/>

## Installation  
### Download `elio`  
```shell
$ git clone https://github.com/cppis/elio
$ cd elio
```

> Now, **$PWD** is the root path.  

<br/>

### [Setting `Skaffold` on Windows](docs/setting.skaffold.md)  
`Skaffold` settings on windows for continuous developing a Kubernetes-native app.  

<br/><br/><br/>

## Run herald  
### using `Skaffold`  
To use the `Skaffold`, you need thd following the [Setup `Skaffold`](#setup-skaffold).  
To run `herald` using `Skaffold`,  
run the following command in the Project root directory:  
```shell
$ skaffold -f app\herald\k8s-resources\skaffold.yaml dev -p local
```

> To change detection triggered to manual mode, use option `--trigger=manual`.  

Or, to run `herald` in debugging mode using `Skaffold`, run the following command:  
```shell
$ skaffold -f app\herald\k8s-resources\skaffold.yaml debug -p local
```

<br/><br/><br/>

## Test  
You can test echo easily by using telnet.  

app protocol is custom `t2p` like http.  
procotol header is separated by newline(`\n` or `\r\n`).  
And packet delimiter is double newline(`\n\n` or `\r\n\r\n`).

### echo: echo message    
  ```
  echo<newline>
  {message}<newline><newline>
  ```
### sub: subcribe to topic    
  ```
  sub<newline>
  {topic}<newline><newline>
  ```
### unsub: unsubcribe from topic  
  ```
  unsub<newline>
  {topic}<newline><newline>
  ```
### pub: publish message to topic  
  ```
  pub<newline>
  {topic}<newline>
  {message}<newline><newline>
  ```

<br/><br/><br/>

## Reference  
* [Building K8S cluster of EMQ X starting from scratch](https://www.emqx.com/en/blog/emqx-mqtt-broker-k8s-cluster)  


<br/><br/><br/>

## TO-DO  
* Run `emqx` as stateful cluster in kubernetes  
