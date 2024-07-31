## tbd

A tiny blog API

### 🚪 user sign up

```yaml
path: /v1/users/create
method: POST
auth: no
body:
    - name
    - email (must be a valid email)
    - password (must be between 8 and 72 characters)
returns:
    - code: 200
        meaning: valid request (json contains id, name and email)
    - code: 400
        meaning: content type is not application/json,body contains badly formed JSON etc
    - code: 422
        meaning: request body contains invalid data (invalid email or password)
    - code: 500
        meaning: request could no longer be processed
```

### 🔓 user log in

```yaml
path: /v1/users/login
method: POST
auth: no
body:
    - email (must be a valid email)
    - password (must be between 8 and 72 characters)
returns:
    - code: 200
        meaning: valid request (json contains jwt for authenticating user)
    - code: 400
        meaning: content type is not application/json,body contains badly formed JSON etc
    - code: 422
        meaning: request body contains invalid data (invalid email or password)
    - code: 403
        meaning: invalid email or password
    - code: 500
        meaning: request could no longer be processed
```

### 🆕 create post

```yaml
path: /v1/posts/create
auth: yes
method: POST
body:
    - title (post title)
    - body (post body)
returns:
    - code: 200
        meaning: valid request (json contains id, time created, author id, title and body)
    - code: 400
        meaning: content type is not application/json,body contains badly formed JSON etc
    - code: 422
        meaning: request body contains invalid data (invalid title or missing body)
    - code: 500
        meaning: request could no longer be processed
```

### 🔍 view post

```yaml
path: /v1/posts/{id}
auth: yes
method: POST
params:
    - id: int
        meaning: id of the post to be viewed
body:
    - title (post title)
    - body (post body)
returns:
    - code: 200
        meaning: valid request (json contains id, time created, author id, title and body)
    - code: 422
        meaning: request body contains invalid data (invalid post id)
    - code: 404
        meaning: post not found
    - code: 500
        meaning: request could no longer be processed
```

### 🔒 auth

```yaml
type: token
kind: bearer token
style: jwt
header: Authorization
lifespan: 2 hours 58 minutes
generator: user login
```

Happy using
