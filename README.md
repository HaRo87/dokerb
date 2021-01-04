# Delphi Planning Poker (Doker) Backend

Just looking for a pet project which can be used by me to learn something about
[Fiber](https://gofiber.io/).

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

3. This projects uses

* [Fiber](https://gofiber.io/)
* [Create Go App](https://create-go.app/)

4. Run project by this command from within the `backend` dir:

```bash
task -s
```

> I am using `Taskfile` as task manager for running the project on a local machine by default. 

## ⚙️ Configuration

```yaml
# ./configs/apiserver.yml

# Server config
server:
  host: 0.0.0.0
  port: 5000

# Database config
database:
  host: 127.0.0.1
  port: 5432
  username: postgres
  password: 1234

# Static files config
static:
  prefix: /
  path: ./static
```

## ⚠️ License

MIT &copy; [HaRo87](https://github.com/HaRo87).


