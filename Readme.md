# Go Git Sync

Tool to sync git repository with other repo


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