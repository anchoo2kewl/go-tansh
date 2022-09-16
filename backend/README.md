# go-tansh
A rewrite of Tansh.us using Go

The easiest way to run is using docker.

For production machines, use the main `Dockerfile`

```
$ docker build -t biswas/tansh:v0.1 .
```

For dev, use another dockerfile, `Dockerfile.dev`

```
$ docker build -f Dockerfile.dev -t biswas/tansh:v0.1-SNAPSHOT .
```

To run the image. From this directory:

```
# First run postgres
$ docker run --name postgres -p 5432:5432 --env-file ../.env -d postgres
# Find IP address of postgres and change the DB host in the .env file

# To run the migrations
$ export $(xargs < ../.env)
$ export POSTGRESQL_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432/$POSTGRES_DB?sslmode=disable"
$ migrate -database ${POSTGRESQL_URL} -path migrations up
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

To backup the database:

# On the remote machine
$ export $(xargs < ../.env)
# A docker restart may be required
$ docker restart `docker container ls  | grep 'db' | awk '{print $1}'`
$ docker exec -t `docker container ls  | grep 'db' | awk '{print $1}'` pg_dumpall -c -U $POSTGRES_USER | gzip > ~/dump_$(date +"%Y-%m-%d_%H_%M_%S").gz

# On the local machine
$ export $(xargs < ../.env)
# Clear data
$ docker-compose down && docker volume rm docker-tansh_db-data
$ gunzip < dump_*.gz | docker exec -i `docker container ls  | grep 'db' | awk '{print $1}'` psql -U $POSTGRES_USER -d $POSTGRES_DB

# Create roles before creating users:

# It is important that the readeer role is the first. Since all users are created with this default role.
$ docker exec -i `docker container ls  | grep 'db' | awk '{print $1}'` psql -U $POSTGRES_USER -d $POSTGRES_DB -c "insert into roles (name, description) values ('viewer','Only has readonly access');"

$ docker exec -i `docker container ls  | grep 'db' | awk '{print $1}'` psql -U $POSTGRES_USER -d $POSTGRES_DB -c "insert into roles (name, description) values ('editor','Can edit data but no admin privileges');"

$ docker exec -i `docker container ls  | grep 'db' | awk '{print $1}'` psql -U $POSTGRES_USER -d $POSTGRES_DB -c "insert into roles (name, description) values ('admin','Is the Superuser of this system');"

# Check if creation has worked fine:
$ docker exec -i `docker container ls  | grep 'db' | awk '{print $1}'` psql -U $POSTGRES_USER -d $POSTGRES_DB -c "select * from roles"