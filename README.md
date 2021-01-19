# go-rest-template

This repository is used as a template for a Go REST API server (mainly for myself).

Provided `Makefile` can compile, test, and create a server and a database (PostgreSQL) containers.

# Storage 

The server can be initialized as an HTTP or HTTPS, depends on the provided configuration.
Additionally, it's possible to set database type: in-memory or external.

## In-memory

All data is stored in server memory.

## External

Data is stored in an external database (currently it's a PostgreSQL database).

# Authorization

Authorization can be set up in 2 modes (session and token).

## Session

This sets the server in server-side session management, meaning that session is stored in 
the server memory and the session ID is set as a `session cookie` in a client.

NB! If server is NOT HTTPS then session cookie will not be sent to the client.

## Token

`WIP`

# TLS crypto material creation:
```
openssl ecparam -name secp384r1 -genkey -noout -out tls.key
openssl req -new -x509 -key tls.key -out tls.crt -days 365
```
