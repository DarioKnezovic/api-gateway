# API Gateway

This repository contains the source code for the API Gateway application.

# Folder Structure

The folder structure of this project is organized as follows:

- `config/`: Contains configuration files and code for handling configuration settings.
    - `config.go`: Handles loading and accessing configuration values.

- `handlers/`: Contains the HTTP request handlers for the API Gateway.
    - `login_handler.go`: Implements the login endpoint.
    - `register_handler.go`: Implements the register endpoint.

- `middleware/`: Contains middleware functions for the API Gateway.
    - `auth_middleware.go`: Implements the authentication middleware.

- `Dockerfile`: Defines the instructions to build the Docker image for the application.
- `Makefile`: Contains make targets for building, running, and hot-reloading the application.


# Docker Image

The Docker image for this application is named `api-gateway`.

## Build

To build the Docker image, run the following command:

```bash
make build
```
This command will build the Docker image using the specified image name (api-gateway).

## Run 

To run a container from the Docker image, use the following command:

```bash
make run
```

This command will start a container from the Docker image api-gateway and expose port 8080 of the container to port 8080 of the host machine.

## Hot reload

To enable hot reloading during development, use the following command:

```bash
make hot-reload
```

This command will start a container from the Docker image api-gateway and use air for hot reloading. It will mount the current directory as a volume inside the container and set the working directory to /app. Any changes made to the code on the host machine will be automatically reflected inside the container.

## Development

To simplify the build and hot reloading process during development, you can use the dev target:

```bash
make dev
```

This command will build the Docker image (make build) and start the hot reloading process (make hot-reload) in one command.

