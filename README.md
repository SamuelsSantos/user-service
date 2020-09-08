## Build

  make

## Run tests

  make test

## API requests 

### Add user

```
curl -X "POST" "http://localhost:8080/v1/user" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "password" : "challenge",
	"email" : "zelda@hash.com.br",
	"first_name" : "zelda",
	"last_name" : "zica",
	"date_of_birth" : "2006-01-02T02:00:00Z"
}'

```
### Search user

```
curl "http://localhost:8080/v1/user?name=ozzy" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Show users

```
curl "http://localhost:8080/v1/user" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```


