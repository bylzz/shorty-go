# Url Shortener built with Go

A simple URL shortener using Go and Mongo.

This project was built using [Echo](https://echo.labstack.com/) and offers the option between two data stores:

* [Badger](https://github.com/dgraph-io/badger) - Embedded Go Key/Value Database as a single standalone executable 
* [Mongo](https://github.com/mongodb/mongo-go-driver) - Docker container that can run as a microservice

In addition, the Mongo database uses [go-cache](github.com/patrickmn/go-cache) as a ttl cache to prevent burdening the database with redundant requests.  

## Running as Docker container

```
make build
make start
```

Service will be available on http://localhost:8080


## Running standalone executable

```
make standalone
make run
```

Service will be available on http://localhost:8080

## Running development

```
make dev
```

## Running tests

```
make test
```

## API

The API is pretty simple.

```
Authentication required - uri?key=:authKey

POST    /shorten                    body{ url:String }
GET     /links?l=:limit&s=:skip     *Mongo only
GET     /links/:code
DELETE  /links/:code

No Authentication required

GET     /               Landing page
GET     /:code          Redirect to long url
GET     /*              404 page
```

## Config and Env variables

This service uses [Viper](https://github.com/spf13/viper) for it's configuration. The config.yaml contains all the 
configurable variables. You can also override any variables as environment variables. You will see examples of this
in the docker-compose.yml. You can also set the variables from the command line.

```
ENV=prod STORE=mongo ./bin/urlshortener
```

Feel free to fork it, hack it and use it any way you please.

**MIT License**