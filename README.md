# go-todo
Classic to-do/scheduler REST API written in Go, featuring graceful shutdown, rate limiter. An example front end implementation can be found here:
http://github.com/pafirmin/calendar-app

# Running the prject
- `git clone https://github.com/pafirmin/go-todo.git`
- `cp .makerc.example .makerc`

You will need to fill out at least the DEV_DB_ADDR and DEV_JWT_SECRET env variables

- `make run/app`

To run the migrations, you will need to install [Migrate](https://github.com/golang-migrate/migrate) and run `make db/migrations/up`

