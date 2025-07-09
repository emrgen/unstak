# unpost

A simple microservice to manage posts.

## Design

This microservice is designed to be a simple CRUD service for managing tinyposts. 
It is built using gRPC and protobuf for communication between the client and the server. 
The server is built using Go and the client is built using Python.

When user account is created, a default unpost is created for the user with the same name as the username.
The user can further add other users to the project as project members.

The user can create token for project which can be used to access the project resources without login.

## Progress

- [x] Create post
- [x] Update post
- [x] List posts
- [x] Delete post
- [ ] Erase post
- [ ] Search post
- [ ] Create tag
- [ ] List tag
- [ ] Delete tag
- [ ] Create tier
- [ ] Get tier
- [ ] List tier
- [ ] Update tier
- [ ] Delete tier
- [ ] List subscriptions
- [ ] list post authors
- [ ] add post author
- [ ] remove post author
- [ ] update post author


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
