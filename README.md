# Delphi Planning Poker (Doker) Backend
![DokerB Logo](backend/static/img/Doker_Logo_DokerB_small.png?raw=true)

![Testing Go Code](https://github.com/HaRo87/dokerb/workflows/Testing%20Go%20Code/badge.svg?branch=main&event=push)
[![codecov](https://codecov.io/gh/HaRo87/dokerb/branch/main/graph/badge.svg?token=YNELZZ65S1)](https://codecov.io/gh/HaRo87/dokerb)

<img src="https://img.shields.io/badge/Go-1.15+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;<img src="https://img.shields.io/badge/license-mit-red?style=for-the-badge&logo=none" alt="license" />

The idea behind Doker is to have a web based service which
allows developers to kinda play [Planning Poker](https://en.wikipedia.org/wiki/Planning_poker)
but using the [Delphi Method](https://en.wikipedia.org/wiki/Delphi_method) for effort estimates.

## ⚡️ Quick start

1. Clone this repository

2. Install the prerequisites 

* [Go](https://golang.org/doc/install)
* [Task](https://taskfile.dev/#/)
* [Swag](https://github.com/swaggo/swag)

3. This projects uses

* [Fiber](https://gofiber.io/)
* [Create Go App](https://create-go.app/)
* [fiber-swagger](https://github.com/arsmn/fiber-swagger)

4. Run project by this command from within the `backend` dir:

```bash
task -s
```

> I am using `Taskfile` as task manager for running the project on a local machine by default. 

Then you should be able to navigate to `http://127.0.0.1:<port>` where `port` is defined by your config
file -> default is `5000`.

This should display the general info page which also provides links to 
further documentation and the Swagger documentation. You can then use
either the interactive Swagger documentation or a tool like 
[HTTPie](https://httpie.io) for trying out the API.

Creating a new session should be as simple as running:

```bash
http POST http://127.0.0.1:5000/api/sessions
```

to create a new session where the response should look something like:

```json
{
    "message": "ok",
    "route": "/sessions/eaf27c59ecdf0db4e165c4f940e176ec"
}
```

Then you can also do things like adding a user:

```bash
http POST http://127.0.0.1:5000/api/sessions/eaf27c59ecdf0db4e165c4f940e176ec/users name="Tigger"
```

which should produce an output like:

```json
{
    "message": "ok",
    "route": "/sessions/eaf27c59ecdf0db4e165c4f940e176ec/users/Tigger"
}
```

## ⚙️ Configuration

```yaml
# ./configs/apiserver.yml

# Server config
server:
  host: 0.0.0.0
  port: 5000

# Database config
database:
  location: my.db

# Static files config
static:
  prefix: /
  path: ./static
```

## Docker Container

In case you want to run DokerB in a Docker container you can use the 
Dockerfile located inside the `backend` directory. 

Just run:

```bash
docker build -t dokerb .
```

## ⚠️ License

MIT &copy; [HaRo87](https://github.com/HaRo87).


