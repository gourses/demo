# Demo Project
Minimal Go project to verify your develepment environment.

## Requirements

### Must have:
- [Make](https://www.gnu.org/software/make/)
- [Go](https://go.dev/) >= 1.13
- Docker or Podman with compose for integration tests with [testcontainers](https://golang.testcontainers.org/)

### Optional:
- awk for `make help`

## Verifying your environment

Makefile should automatically install required Go version and tools besides the ones listed in requirements.

Running:
```sh
$ make
```
Should result into:
```sh
Usage:
  make <recipe>

Recipes:
  help                    Show help
  build                   Build binary
  run                     Execute binary
  unit-test               Run unit tests
  integration-test        Run integration tests
  merge-coverages         Merge unit and integration coverages
  show-coverage           Open coverage report in browser
  clean                   Remove all build and test artifacts
  clean-tools             Remove all tools
```

Running:
```sh
$ make build
```
Should result into:
```sh
# Some Go installation log if this is the first run.
.tools/go/go1.20.3 build -o target/bin/demo ./cmd/demo
Binary target/bin/demo built successfully!
```

Running:
```sh
$ make lint
```

Should result into:
```sh
# Some golanci-lint installation log if this is the first run.
.tools/github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2 run
api/api.go:38:9: error returned from external package is unwrapped: sig: func (*net/http.Server).ListenAndServe() error (wrapcheck)
	return err
	      ^
app.go:32:9: error returned from external package is unwrapped: sig: func (*github.com/gourses/demo/api.API).Run() error (wrapcheck)
	return a.Run()
	      ^
make: *** [Makefile:53: lint] Error 1
```

Running:
```sh
$ make image
```

Should result into:
```sh
CGO_ENABLED=0 .tools/go/go1.20.3 build -o target/bin/demo ./cmd/demo
Binary target/bin/demo built successfully!
docker build -t demo:latest .
Emulate Docker CLI using podman. Create /usr/etc/containers/nodocker to quiet msg.
STEP 1/3: FROM scratch
STEP 2/3: COPY /target/bin/demo /
--> 053db692c185
STEP 3/3: ENTRYPOINT ["/demo"]
COMMIT demo:latest
--> c0b8bb22e916
Successfully tagged localhost/demo:latest
c0b8bb22e916f9f4d8b0e861e9aef2c421d717722aeaae083c43dfb2df35a1c1
```

Running:
```sh
$ make show-coverage
```

Should result into:
```sh
.tools/go/go1.20.3 test -race -covermode atomic -coverprofile=target/coverage/unit.out ./...
?   	github.com/gourses/demo	[no test files]
?   	github.com/gourses/demo/api	[no test files]
?   	github.com/gourses/demo/cmd/demo	[no test files]
ok  	github.com/gourses/demo/storage	0.017s	coverage: 46.2% of statements
.tools/go/go1.20.3 test -race -covermode atomic -coverprofile=target/coverage/integration.out -tags=integration -coverpkg=github.com/gourses/demo,github.com/gourses/demo/api,github.com/gourses/demo/cmd/demo,github.com/gourses/demo/storage ./cmd/demo
ok  	github.com/gourses/demo/cmd/demo	6.940s	coverage: 90.9% of statements in github.com/gourses/demo, github.com/gourses/demo/api, github.com/gourses/demo/cmd/demo, github.com/gourses/demo/storage
.tools/github.com/wadey/gocovmerge@master target/coverage/unit.out target/coverage/integration.out > target/coverage/merged.out
.tools/go/go1.20.3 tool cover -html=target/coverage/merged.out
```
and coverage report should open in browser.
