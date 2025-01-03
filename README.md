# unpost

A simple microservice to manage tinyposts.

## Design

This microservice is designed to be a simple CRUD service for managing tinyposts. 
It is built using gRPC and protobuf for communication between the client and the server. 
The server is built using Go and the client is built using Python.

When user account is created, a default unpost is created for the user with the same name as the username.
The user can further add other users to the project as project members.

The user can create token for project which can be used to access the project resources without login.

## Features

- [x] Create project
- [x] Update project
- [x] List project
- [x] Create space
- [x] Update space
- [x] List space
- [x] Create unpost
- [ ] Update unpost
- [x] Delete unpost
- [ ] Erase unpost
- [x] Get unpost
- [x] List unpost
- [ ] Search unpost
- [x] Document unpost
- [x] Create tag (tags which are not used for more than 10 days will be deleted automatically)
- [x] List tag
- [ ] Delete tag

## Installation

```bash
# install initial dependencies(it will fail but that's fine, still need to run it)
make deps
# build proto
make protoc
# install all dependencies
make deps
```

## Usage
