# D7024E Kademlia_lab

# Setup
## Install dependencies
The following dependencies is required to run this program:
- Docker
- Docker-compose
- Golang


## Build docker image
To build the docker image, run the following commands from the root directory of this project:
```
docker build . -t kadlab
```

## Run containers with Docker Compose
Modify the number of nodes in the docker-compose file.
```
replicas: 50
```
Run this command from the root directory of the project to start the containers:
```
docker-compose up --build
```

To terminate the containers, run command:
```
docker-compose down
```

# CLI APP
The CLI app has three main functions.
Available commands:
- put \<string> (stores data on k closest nodes to hash)	
- get \<hash> (fetches data object with this hash if it is stored in the network)
- exit (terminates this node)

From the container terminal, run the following command to start the CLI app:
```
/cli
```