# go-tansh
A rewrite of Tansh.us using Go

The easiest way to run is using docker.

For production machines, use the main `Dockerfile`

```
$ docker build -t biswas/tansh:v0.1 .
```

For dev, use another dockerfile, `Dockerfile.dev`

```
$ docker build -f Dockerfile.production -t biswas/tansh:v0.1-SNAPSHOT .
```

To run the image:

```
$ docker run -d --env-file ../.env --name tansh-prod -it --rm -p 3001:3000 -v $PWD:/go/src/tansh biswas/tansh:v0.1
# or to run the dev image:
$ docker run -d --env-file ../.env --name tansh -it --rm -p 3000:3000 -v $PWD:/go/src/tansh biswas/tansh:v0.1-SNAPSHOT
```

The difference between dev and prod are two main things:
* Prod is compiled and hard to change the view any changes without re-building. Dev allows changes without restarting the container.
* The Prod image is about 9 Mb. The Dev image is about 1200 Mb.

To run without docker, just pull the mods and run:


```
$ go mod download github.com/go-chi/chi/v5
$ go mod download   docgen
$ go mod download github.com/go-chi/render
$ go run main.go
```