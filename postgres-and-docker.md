## Installing Docker

```bash
# Make sure docker service is up and running
# Check all services
sudo service --status-all
# Start docker service
sudo service docker start
```

> Create docker-compose.yml

```yml
version: "3.9"

services:
  # Our Postgres database
  db: # The service will be named db.
    image: postgres # The postgres image will be used
    restart: always # Always try to restart if this stops running
    environment: # Provide environment variables
      POSTGRES_USER: baloo # POSTGRES_USER env var w/ value baloo
      POSTGRES_PASSWORD: junglebook
      POSTGRES_DB: lenslocked
    ports: # Expose ports so that apps not running via docker compose can connect to them.
      - 5432:5432 # format here is "port on our machine":"port on container"

  # Adminer provides a nice little web UI to connect to databases
  adminer:
    image: adminer
    restart: always
    environment:
      ADMINER_DESIGN: dracula # Pick a theme - https://github.com/vrana/adminer/tree/master/designs
    ports:
      - 3333:8080
```

```bash
docker compose up
```
