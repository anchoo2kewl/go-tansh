# Docker-tansh
Run the various Tansh dependencies with docker-compose:

To run in development:

```
docker-compose up
```

This would allow files to be quickly changed without rebuilding anything. However, the size of the container is big. For production, run this:

```
docker-compose -f docker-compose.production.yml up
```