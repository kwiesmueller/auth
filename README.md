# Auth

## Install

`go get -u github.com/bborbe/auth/bin/auth_server`

`go get -u github.com/siddontang/ledisdb/cmd/ledis-server`

## Run

Start ledis database

```
ledis-server \
-databases=1 \
-addr=localhost:5555
```

Start auth-server

```
auth_server \
-logtostderr \
-v=2 \
-port=6666 \
-ledisdb-address=localhost:5555 \
-auth-application-password=test123 \
-prefix=
```

## Authorize Header

`echo -n 'auth:test123' | base64`

`Authorization: Bearer YXV0aDp0ZXN0MTIz` 

## Actions Global

### Health-Check

```
curl \
-X GET \ 
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/healthz
```

### Readiness-Check

```
curl \
-X GET \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/readiness
```

## Actions Api 1.0

### Register User

`echo -n 'tester:secret' | base64`

```
curl \
-X POST \
-d '{ "authToken":"dGVzdGVyOnNlY3JldA==","user":"tester" }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user
```

### Unregister User

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/token/dGVzdGVyOnNlY3JldA==
```

### Delete User

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user/tester
```

### Verify Login

```
curl \
-X POST \
-d '{ "authToken": "abc", "groups": ["admin"] }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" http://localhost:6666/api/1.0/login
```

### Create application

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" http://localhost:6666/api/1.0/application
```

### Delete application

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/application/1337
```

### Get application

```
curl \
-X Get \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/application/1337
```

### Add authoken to existing user 

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/application
```

### Remove authtoken from existing user 

```
curl \
-X DELETE \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/application
```

### Add authoken to existing user 

```
curl \
-X POST \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user_group
```

### Remove authtoken from existing user 

```
curl \
-X DELETE \
-d '{ ... }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user_group
```

### Create user data

```
curl \
-X POST \
-d '{ "keya":"valuea", "keyb":"valueb" }' \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user/tester/data
```

### Get user data

```
curl \
-X GET \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user/tester/data
```

### Get user data key

```
curl \
-X GET \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user/tester/data/keya
```

### Delete user data

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user/tester/data
```

### Delete user data key

```
curl \
-X DELETE \
-H "Authorization: Bearer YXV0aDp0ZXN0MTIz" \
http://localhost:6666/api/1.0/user/tester/data/keya
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
