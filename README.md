# Content building block

Go project to provide rest service for rokwire building block content results.

## Set Up

### Prerequisites

MongoDB v4.2.2+

Go v1.16+

### Environment variables
The following Environment variables are supported. The service will not start unless those marked as Required are supplied.

Name|Value|Required|Description
---|---|---|---
CONTENT_PORT | < value > | yes | The port number of the listening port
CONTENT_HOST | < value > | yes | Host name
CONTENT_MONGO_AUTH | <mongodb://USER:PASSWORD@HOST:PORT/DATABASE NAME> | yes | MongoDB authentication string. The user must have read/write privileges.
CONTENT_MONGO_DATABASE | < value > | yes | MongoDB database name
CONTENT_MONGO_TIMEOUT | < value > | no | MongoDB timeout in milliseconds. Set default value(500 milliseconds) if omitted
CONTENT_OIDC_CLIENT_IDS | < value > | yes | OIDC admin client id
CONTENT_OIDC_PROVIDER | < value > | yes | OIDC provider
CONTENT_PHONE_SECRET | < value > | yes | Phone secret
CORE_BB_HOST | < value > | yes | Core BB host url
ROKWIRE_API_KEYS | <value1,value2,value3> | yes | Comma separated list of rokwire api keys
AWS_ACCESS_KEY_ID | < value > | yes | AWS Access key ID
AWS_SECRET_ACCESS_KEY | < value > | yes | AWS Secret access ket
S3_BUCKET | < value > | yes | AWS S3 bucket name
S3_REGION | < value > | yes | AWS S3 region name
TWITTER_FEED_URL | < value > | yes | Twitter Feed base URL
TWITTER_ACCESS_TOKEN | < value > | yes | Twitter Bearer access token
DEFAULT_CACHE_EXPIRATION_SECONDS | < value > | false | Default cache expiration time in seconds. Default: 120
### Run Application

#### Run locally without Docker

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Make the project  
```
$ make
...
▶ building executable(s)… 1.9.0 2020-08-13T10:00:00+0300
```

4. Run the executable
```
$ ./bin/content
```

#### Run locally as Docker container

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Create Docker image  
```
docker build -t content .
```
4. Run as Docker container
```
docker-compose up
```

#### Tools

##### Run tests
```
$ make tests
```

##### Run code coverage tests
```
$ make cover
```

##### Run golint
```
$ make lint
```

##### Run gofmt to check formatting on all source files
```
$ make checkfmt
```

##### Run gofmt to fix formatting on all source files
```
$ make fixfmt
```

##### Cleanup everything
```
$ make clean
```

##### Run help
```
$ make help
```

##### Generate Swagger docs
```
$ make swagger
```

### Test Application APIs

Verify the service is running as calling the get version API.

#### Call get version API

curl -X GET -i http://localhost/content/version

Response
```
1.9.0
```

## Documentation

The documentation is placed here - https://api-dev.rokwire.illinois.edu/docs/

Alternatively the documentation is served by the service on the following url - https://api-dev.rokwire.illinois.edu/content/doc/ui/
