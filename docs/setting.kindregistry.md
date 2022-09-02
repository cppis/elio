# [Kind 로컬 레지스트리 만들기](https://kind.sigs.k8s.io/docs/user/local-registry/)

쿠버네티스 로컬 이미지 레지스트리가 없는 경우 새로 생성해야 합니다.  

먼저, 다음 명령으로 로컬 도커 레지스트리가 이미 있는 지 확인합니다:  

```bash
# kind-registry 가 실행 중인지 확인합니다.(True: 실행 중)
docker inspect -f '{{.State.Running}}' kind-registry
```

위 명령의 실행 결과가 `True` 가 아니라면, 다음 명령으로 로컬 도커 레지스트리를 생성합니다:  

```bash
# kind-registry 가 없다면 새로 생성합니다.
docker run -d --restart=always -p 127.0.0.1:5001:5000 --name kind-registry registry:2
```

설치가 끝났다면, 다음 명령으로 도커 레지스트리를 확인할 수 있습니다:

```bash
curl http://127.0.0.1:5001/v2/_catalog
  {"repositories":[]}
```

> 자세한 내용은 [Kind - 로컬 레지스트리](https://kind.sigs.k8s.io/docs/user/local-registry/) 를 참고하세요.  

> 만약, Kind 클러스터에서 로컬 Docker 레지스트리에 접근하지 못한다면,  
> 다음 명령을 실행하여 Kind 클러스터의 네트워크를 로컬 Docker 레지스트리의 네트워크와 연결합니다:  
> ```bash
> docker network connect "kind" "kind-registry"
> ```
