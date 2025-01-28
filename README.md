# Let-s-Go-snippetbox

Snippetbox share the same concept of [GitHub Gist](https://gist.github.com/) where people can share snippets of text.

# Key Takeaways

- HTML Templating
- RESTful Routing
- Configuration Management
- Level logging and centralized error handling
- MySQL Database
- Middlewares
- Session Management
- Security (HTTPS,OWASP Secure Heards and CSRF)
- Authentication and Authorization
- Contexts
- Testing
- Change Password

# How to Use

1. Clone the project repository to your local machine:

   ```bash
     git clone git@github.com:GergesHany/Let-s-Go-snippetbox.git
   ```

2. Run the project:

   ```bash
     go run ./cmd/web
   ```

3. Open your browser and navigate to `http://localhost:4000/`.

## Requirements

- Go 1.18
- MySQL

## Configuration Management

- Snippetbox supports configuration management for the following flags:

  - `addr` - Port Address (default is ':4000')
  - `dsn` - DSN (default is 'web:pass@/snippetbox?parseTime=true')\*
  - `debug` - Debug mode (default is false)

- It is very important to set 'parseTime=true', otherwise the mapping for dates may not work properly

### Example:

```
$ go run ./cmd/web -debug -port=:8080 -dsn=any:pass@/mybox?parseTime=true
```

It will run on _port_ 8080 with _dsn_ 'any:pass@/mybox?parseTime=true' and _debug_ mode enabled.

<hr>

## Routes

| HTTP Method | Route                      | Middleware/Handler                                  | Notes                           |
| ----------- | -------------------------- | --------------------------------------------------- | ------------------------------- |
| GET         | `/static/*filepath`        | `http.FileServer(http.FS(ui.Files))`                | Serves static files.            |
| GET         | `/ping`                    | `ping`                                              | Health check endpoint.          |
| GET         | `/`                        | `dynamic.ThenFunc(app.home)`                        | Unprotected dynamic route.      |
| GET         | `/about`                   | `dynamic.ThenFunc(app.about)`                       | Unprotected dynamic route.      |
| GET         | `/snippet/view/:id`        | `dynamic.ThenFunc(app.snippetView)`                 | Unprotected dynamic route.      |
| GET         | `/user/signup`             | `dynamic.ThenFunc(app.userSignup)`                  | Unprotected dynamic route.      |
| POST        | `/user/signup`             | `dynamic.ThenFunc(app.userSignupPost)`              | Unprotected dynamic route.      |
| GET         | `/user/login`              | `dynamic.ThenFunc(app.userLogin)`                   | Unprotected dynamic route.      |
| POST        | `/user/login`              | `dynamic.ThenFunc(app.userLoginPost)`               | Unprotected dynamic route.      |
| GET         | `/snippet/create`          | `protected.ThenFunc(app.snippetCreate)`             | Protected (authenticated-only). |
| POST        | `/snippet/create`          | `protected.ThenFunc(app.snippetCreatePost)`         | Protected (authenticated-only). |
| POST        | `/user/logout`             | `protected.ThenFunc(app.userLogoutPost)`            | Protected (authenticated-only). |
| GET         | `/account/view`            | `protected.ThenFunc(app.accountView)`               | Protected (authenticated-only). |
| GET         | `/account/password/update` | `protected.ThenFunc(app.accountPasswordUpdate)`     | Protected (authenticated-only). |
| POST        | `/account/password/update` | `protected.ThenFunc(app.accountPasswordUpdatePost)` | Protected (authenticated-only). |
