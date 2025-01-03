# Design

## Create Post

When the create post request is made, the following steps are taken:

1. The request is received by the server.
2. The server validates the request.
3. The server check if the user is authenticated.
4. The server injects user space permissions into the request.
5. The server creates a new post.
   1. The server creates a new document in the document store.
   2. The server creates a new post in the unpost store with the document id, user id, and space id.
6. The server returns the post.

## Get Post

When the get post request is made, the following steps are taken:

1. The request is received by the server.
2. The server validates the request.
3. The server check if the user is authenticated.
4. The server injects user space permissions into the request.
5. The server retrieves the post from the unpost store.
6. The server retrieves the document from the document store using the document id.
7. The server returns the post.