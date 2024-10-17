# bookshelf

## Overview

This project is a simple RESTful API built using Golang for managing books and authors, with user authentication and authorization using JWT tokens. It supports CRUD (Create, Read, Update, Delete) operations for both books and authors, along with user authentication features. The API ensures hash validation on each request to maintain data integrity, and the entire setup is Dockerized using Docker Compose.

## API Endpoints

### Books

- `GET /books`: Retrieve a list of all books.
- `GET /books/{id}`: Retrieve details of a book by its ID.
- `POST /books`: Create a new book.
- `PUT /books/{id}`: Update an existing book by ID.
- `DELETE /books/{id}`: Delete a book by ID.

### Authors

- `GET /authors`: Retrieve a list of all authors.
- `GET /authors/{id}`: Retrieve an author's details by ID.
- `POST /authors`: Create a new author.
- `PUT /authors/{id}`: Update an existing author by ID.
- `DELETE /authors/{id}`: Delete an author by ID.

### Users

- `POST /auth/register`: Register a new user.
- `POST /auth/login`: Authenticate a user and return a JWT token.

## Prerequisites

- Clone the repository:

  ```bash
  git clone https://github.com/mnaufalhilmym/bookshelf.git
  cd bookshelf
  ```

- Install the dependencies:

  ```bash
  go mod tidy
  ```

## Running The Application

- Ensure that you have created and configured `config.yml`. An example configuration file can be found in `config-example.yml`.

- Run the application

  ```bash
  go run ./cmd
  ```

## Build and Running with Docker

- Build the Docker image:

  ```bash
  docker build . -f Dockerfile -t docker.io/mnaufalhilmym/bookshelf
  ```

- Run the Docker container:

  ```bash
  docker run -p 8080:8080 --name bookshelf -v $(pwd)/config.yml:/config.yml docker.io/mnaufalhilmym/bookshelf
  ```

- Run with Docker Compose:

  ```bash
  docker-compose -f compose.yml up
  ```

## API Documentation

[Postman API Documentation](bit.ly/bookshelf-api-docs)
