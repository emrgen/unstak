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
- [x] Get post
- [x] List post
- [ ] Search post
- [x] Create tag
- [x] List tag
- [x] Delete tag
- [x] Create collection
- [x] List collection
- [x] Delete collection
- [x] Create tier
- [x] Get tier
- [ ] List tier
- [ ] Update tier
- [ ] Delete tier
- [ ] List subscriptions
- [x] Create space
- [x] List space
- [x] Update space
- [ ] delete space
- [ ] Get space
- [x] list space members
- [x] add space member
- [ ] remove space member
- [ ] update space member
- [ ] create offline token with restricted permission to one space
- [ ] list offline token
- [ ] remove offline token
- [ ] delete post
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
