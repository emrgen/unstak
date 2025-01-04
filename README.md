# unpost

A simple microservice to manage posts.

## Design

This microservice is designed to be a simple CRUD service for managing tinyposts. 
It is built using gRPC and protobuf for communication between the client and the server. 
The server is built using Go and the client is built using Python.

When user account is created, a default unpost is created for the user with the same name as the username.
The user can further add other users to the project as project members.

The user can create token for project which can be used to access the project resources without login.

## Features

- [x] Create post
- [ ] Update post
- [x] Delete post
- [ ] Erase post
- [x] Get post
- [x] List post
- [ ] Search post
- [x] Create tag
- [x] List tag
- [ ] Delete tag
- [ ] Create collection
- [ ] List collection
- [ ] Delete collection
- [ ] Create subscription
- [ ] List subscription
- [ ] Delete subscription
- [ ] Update subscription

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
