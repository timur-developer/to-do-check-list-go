# ToDoList-Go

**ToDoList-Go** is a robust and scalable to-do list application built with a microservice architecture. It allows users to create, track, delete, and mark tasks as completed. This project demonstrates key concepts of microservices, event-driven architecture, and efficient task management.

## Features

- **API Service**: Handles HTTP requests for task operations (create, retrieve, delete, update).
- **Database Service**: Manages interactions with a PostgreSQL database, storing task data and managing task statuses. 
- **Optional Kafka Service**: Logs events and task actions for improved observability and tracking.
- **Docker Compose**: Orchestrates the microservices, making it easy to deploy and scale the application.
- **Comprehensive Logging**: Logs incoming requests, outgoing responses, and errors for troubleshooting.
- **Unit and Integration Tests**: Ensures the reliability and correctness of the system. 

## Architecture

The application is built using the following microservices:

- **API Service**: Responsible for accepting HTTP requests from the user and interacting with the database service.
- **Database Service**: Connects to PostgreSQL and handles CRUD operations for tasks.
- **Kafka Service** (optional): Listens for events from the API service and logs actions in Kafka for event tracking.

## Getting Started

To get started with this project locally, follow these steps:

### Prerequisites

- Docker and Docker Compose must be installed on your machine.

### Installation

1. Clone this repository:
    ```bash
    git clone https://github.com/timur-developer/to-do-list-go.git
    ```

2. Navigate to the project directory:
    ```bash
    cd to-do-list-go
    ```

3. Build and run the application using Docker Compose:
    ```bash
    docker-compose up --build
    ```

4. The services will be up and running, and you can start interacting with the API.

## Endpoints

The following API endpoints are available:

- **POST /create**: Create a new task.
- **GET /list**: Retrieve a list of all tasks.
- **DELETE /delete**: Delete a task by its unique ID.
- **PUT /done**: Mark a task as completed by its unique ID. 


Contributions are welcome! Feel free to open an issue or submit a pull request if you have any improvements or suggestions for the project.💻

## License📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.🎉
