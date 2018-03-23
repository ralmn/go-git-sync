# Go Git Sync

Tool to sync git repository with other repo

Developed by [Mathieu "ralmn" HIREL](https://github.com/ralmn) ([@ralmn45](https://twitter.com/ralmn45))

## Build 

### Docker

`docker build --build-args APP_VERSION=$(cat version) .`

### Manual build

go build -ldflags "-X main.version=$(cat version)" -a -installsuffix cgo -o go-git-sync . 

## Usage

### Running server 
The server listen on port *8080*

Run server : 

    ./go-git-sync

### Running with docker 

docker-compose.yml example : 

```
version: '3'
services: 
  go-git-sync:
    image: ralmn/go-git-sync
    ports: 
      -  "127.0.0.1:9246:8080"
    volumes:
      - "/srv/go-git-sync/config.toml:/root/config.toml"
      - "/srv/go-git-sync/id_rsa:/root/.ssh/id_rsa"
      - "/srv/go-git-sync/known_hosts:/etc/ssh/ssh_known_hosts" 
```

### Setup webhook 

For sync repository you need to call a http webhook :  

    http://localhost:8080/webhook/push/<repository name>/<mirror name>?secret=<secret>
    
The secret token is generate on the first start and saved in *config.toml*
    
For exemple with the repository *go-git-sync* and the remote/mirror *gitlab* : 

    http://localhost:8080/webhook/push/go-git-sync/gitlab?secret=********
     

## Config exemple

`config.toml`

```
[[repository]]
    name = "go-git-sync"
    [[repository.mirror]]
        name = "github"
        remote = "git@github.com:ralmn/go-git-sync.git"
    [[repository.mirror]]
        name = "gitlab"
        remote = "2nd remote"


``