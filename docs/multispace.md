# Multi space support

This feature is available in the following versions:
v1.0.0 and later

## Overview

This feature allows you to create multiple spaces in your application. Each space is a separate instance of the application with its own data and configuration. 
This is useful for creating different environments for different users. Users can also create private spaces for their own use. Private space optionally support user pools.

## Space types

There are two types of spaces:

1. Public space: This is the default space type. It is accessible to all users.
2. Private space: This space is accessible only to the user who created it. It can optionally support user pools.

Public spaces are created by user and other users can access them. Posts published in public spaces are visible to all users. Space admins can 
Private spaces are created by the user and are accessible only to the user who created them.

## User pools

User pools are optional for private spaces. When a user creates a private space with user pool, the space is accessible only to the user who created it. 
Other users can access the space only if they create an account in the user pool associated with the space.

## Space CLI commands

### How to create a space

```bash
unpost space create -n <space-name>
```

### How to list spaces

```bash
unpost space list
```

### How to switch to a space

```bash
unpost space switch -n <space-name>
```

### How to delete a space

```bash
unpost space delete -n <space-name>
```

### How to create a private space

```bash
unpost space create -n <space-name> --private
```

### How to create a private space with user pool

```bash
unpost space create -n <space-name> --private --user-pool
```
