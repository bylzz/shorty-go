# Url Shortener built with Go

A simple URL shortener using Go and Mongo.

This project was built using [Echo](https://echo.labstack.com/) and offers the option between two data stores:

* [Badger](https://github.com/dgraph-io/badger) - Embedded Go Key/Value Database as a single standalone executable. Great for light / mid-tier usage.
* [Mongo](https://github.com/mongodb/mongo-go-driver) + [go-cache](github.com/patrickmn/go-cache) - Allows for greater control of data storage and clustered environments. 

Microservice runs on http://localhost:8080 by default.

### Running standalone executable

Give the code a quick spin by building a single exec with no external dependencies. Optional config.yaml file can be used for configuration or you can configure it as part of the command line.

```
make standalone
make run
```

### Running as Docker container

```
make build
make start
```

### Running development

```
make dev
```

### Running tests

```
make test
```

### API

The API is pretty simple.

```
Authentication required - uri?key=:authKey

POST    /shorten                    body{ url:String }
GET     /links?l=:limit&s=:skip     (Mongo only)
GET     /links/:code
DELETE  /links/:code

No Authentication required

GET     /               Landing page
GET     /:code          Redirect to long url
GET     /*              404 page
```

### Config and Env variables

URL shortener uses [Viper](https://github.com/spf13/viper) to handle configuration. The `config.yaml` contains all the 
configurable variables and default values. You can also override any variables as environment variables. You will see examples of this
in the `docker-compose.yml`. You can also set the variables from the command line.

```
ENV=prod STORE=mongo ./bin/urlshortener
```

Feel free to fork it, hack it and use it any way you please.

**MIT License**