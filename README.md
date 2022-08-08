# Docker-tansh
Run the various Tansh dependencies with docker-compose:

To run in production:

```
# Copy the sample
$ cp docker-compose.sample.yml docker-compose.yml
# Copy the .env file
$ cp .env.sample .env
# Replace the PG URL
$ export $(xargs < .env)
$ export POSTGRESQL_URL="postgres:\/\/$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432\/$POSTGRES_DB?sslmode=disable"
# For Mac
$ sed -i '' "s/{{PG_STRING}}/$POSTGRESQL_URL/g" ../docker-compose.yml
$ For Linux
$ sed -i "s/{{PG_STRING}}/$POSTGRESQL_URL/g" ../docker-compose.yml
$ docker-compose up
```

For development, we have a config that allows files to be quickly changed without rebuilding anything. However, the size of the container is big.

```
$ docker-compose -f docker-compose.dev.yml up
```