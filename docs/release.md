## Building consul-debug-read go release

* Set versioning at `internal/read/version.go`

* Set GH Release Token
  * `export GITHUB_TOKEN=$CONSUL_DEBUG_GH_TOKEN`

* Create tag, and push to GitHub: 
  * `git tag -a v1.1.8 -m "Patch release v1.1.8"`
  * `git push origin v1.1.8`

* Release:
  * `goreleaser release --clean`



## Recreate tag post-correction

```shell
git tag -d v1.1.8
git push --delete origin v1.1.8
```


re-run the above