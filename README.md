# Auth

## Install

`go get github.com/bborbe/auth/bin/auth_server`

## Run

Start ledis database

```
ledis-server \
-databases=1 \
-addr=localhost:6380
```

Start auth-server

```
auth_server \
-loglevel debug \
-ledisdb-address localhost:6380 \
-auth-application-password test123
-port 8080
```

## Authorize

`echo -n 'auth:test123' | base64`

`Authorization: Bearer YXV0aDp0ZXN0MTIz` 

## Actions

### Health-Check

```
curl \
-X GET \ 
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/healthz
```

### Readiness-Check

```
curl \
-X GET \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/readiness
```

### Register User

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user
```

### Unregister User

```
curl \
-X DELETE \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user/1337
```

### Verify Login

```
curl \
-X POST \
-d '{ "applicatonName": "test", "applicatonPassword": "test", "connectorName": "test", "connectorUserIdentifier": "test", }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" http://localhost:8080/login
```

### Create application

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" http://localhost:8080/application
```

### Delete application

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/application/1337
```

### Get application

```
curl \
-X Get \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/application/1337
```

### Add authoken to existing user 

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/application
```

### Remove authtoken from existing user 

```
curl \
-X DELETE \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/application
```

### Add authoken to existing user 

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user_group
```

### Remove authtoken from existing user 

```
curl \
-X DELETE \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user_group
```

### Create user data

```
curl \
-X POST \
-d '{ "keya":"valuea", "keyb":"valueb" }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user/tester/data
```

### Get user data

```
curl \
-X GET \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user/tester/data
```

### Get user data key

```
curl \
-X GET \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user/tester/data/keya
```

### Delete user data

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user/tester/data
```

### Delete user data key

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:8080/user/tester/data/keya
```

## Continuous integration

[Jenkins](https://www.benjamin-borbe.de/jenkins/job/Go-Auth/)

## Copyright and license

    Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>
    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are
    met:
    
       * Redistributions of source code must retain the above copyright
         notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
         copyright notice, this list of conditions and the following
         disclaimer in the documentation and/or other materials provided
         with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
