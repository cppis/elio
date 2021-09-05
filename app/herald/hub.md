# cppis/herald  
`herald` means a messanger.  
It is a sample app to test pub/sub of backing `emqx` message broker  
using [`elio`](https://github.com/cppis/elio) library.  

<figure>
  <div style="text-align:center">
    <img src="https://drive.google.com/uc?export=view&id=1kysK0cWuaM-dvw1eeLgB5U6NmmI70ZpW" style="width: 500; max-width: 100%; height: auto" title="nginx-php" />
  </div>
</figure>


<br/><br/>

## Run  
```shell
$ docker run -e HERALD_IN_URL=0.0.0.0:7000 -e HERALD_MQTT_URL=localhost:1883 -p 7000:7000 cppis/herald
```
* HERALD_IN_URL: `herald` listen URL  
* HERALD_MQTT_URL: `herald` backing `mqtt` URL  

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
  If receive messages from subscription, print to *stdout*. 
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
