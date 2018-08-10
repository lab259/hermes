
[![CircleCI](https://circleci.com/gh/lab259/http.svg?style=shield)](https://circleci.com/gh/lab259/http)
[![codecov](https://codecov.io/gh/lab259/http/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/http)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/http)](https://goreportcard.com/report/github.com/lab259/http)

# HTTP

This is one more web development framework for GO using the fasthttp as base.
Or, maybe, it is just a set of utility of helpful functions that apply some
sugar to the fasthttp "daily" usage.

# Extra Features

* [Grouping](#grouping)
* [Middlewares](#middlewares)
* [Thin JSON Layer](#thin-json-layer)

## Routing

`Routable` is described as an interface which have the DELETE, GET, HEAD,
OPTIONS, PATCH, POST, PUT methods for routing. This aims to replace the
fasthttprouter implementation ([#fasthttprouter](check why here)).

In order to get things going, the below example shows how to define an endpoint
serving a GET method:

```go
	router := http.NewRouter()

	router.GET("/api/v1/user", func(ctx *Context) {
		fmt.Fprint(ctx, "the user is: Snake Eyes")
	})
```

### Grouping

When dealing with routes, groups are awesome!

```go
	router := http.NewRouter()

	apiGroup := router.Group("/api") // This retuns a `Routable` that can be used
	                                 // to create other subgroups or define routes.

	apiv1 := apiGroup.Group("/v1")   // This is what we define as a subgroup.

	apiv1.GET(                       // Now a definition of the route itself.
		"/user",
		func(ctx *Context) {
			fmt.Fprint(ctx, "the user is: Snake Eyes")
		},
	)
	// many other routes using `apiv1` ...
```

There is no difference from using, or not, grouping into the routes definition.
Internally, the implementation ends up joining all the routes of the group with
the group prefix. Hence, groups will not affect performance.

### Middlewares

In order to provide a more flexible API, Middleware supports were added.

Middlewares can implement some verification or extension logic and decide
whether or not continue to run the "next" middleware/handler.

```go
	router := http.NewRouter()

	apiGroup := router.Group("/api")
	apiv1 := apiGroup.Group("/v1")
	apiv1.GET(
		"/user",
		func(ctx *Context) {
			fmt.Fprint(ctx, "the user is: Snake Eyes")
		},
		func(ctx *Context, next Handler) {
			// This is a middleware that could do something smart...
			next(ctx)
		},
	)
```

Yet, you can also define multiple middlewares for each route and their priority
will be from the left to the right.

```go
	router := http.NewRouter()

	apiGroup := router.Group("/api")
	apiv1 := apiGroup.Group("/v1")
	apiv1.GET(
		"/user",
		func(ctx *Context) {
			fmt.Fprint(ctx, "the user is: Snake Eyes")
		},
		func(ctx *Context, next Handler) {
			// This is a middleware that could do something smart...
			next(ctx)
		},
		func(ctx *Context, next Handler) {
			// This is a second middleware for the endpoint...
			next(ctx)
		},
	)
```

Middlewares are also supported on groups:

```go
	router := http.NewRouter()

	apiGroup := router.Group("/api", func(ctx *Context, next Handler) {
		// This is a middleware that could do something smart...
		next(ctx)
	})
	apiv1 := apiGroup.Group("/v1", func(ctx *Context, next Handler) {
		// Yet another middleware applied just for this subgroup...
		next(ctx)
	})
	apiv1.GET(
		"/user",
		func(ctx *Context) {
			fmt.Fprint(ctx, "the user is: Snake Eyes")
		},
		func(ctx *Context, next Handler) {
			// This is a middleware that is applied just for this endpoint
			next(ctx)
		},
	)
```

Again, there is no performance difference when using middlewares in a specific
route or in a whole group. The internal implementation will append both
middleware definitions into one big sequence of middlewares for each route.

### Thin JSON layer

For simple sake of ease the use of sending and receiving JSON objects `SendJson`
and `BodyJson` methods were added to the `Context`.

#### Sending a JSON

Above an example of sending a JSON document:

```go
	router := http.NewRouter()

	apiGroup := router.Group("/api")
	apiv1 := apiGroup.Group("/v1")
	apiv1.GET(
		"/user",
		func(ctx *Context) {
			ctx.SendJson(map[string]interface{}{
				"name": "Snake Eyes",
				"email": "s.eyes@gijoe.com",
			})
		},
	)
```

### Receiving a JSON

Above an example of sending a JSON document:

```go
	router := http.NewRouter()

	apiGroup := router.Group("/api")
	apiv1 := apiGroup.Group("/v1")
	apiv1.POST(
		"/user",
		func(ctx *Context) {
			user := make(map[string]interface{})
			err = ctx.BodyJson(&user)
			if err != nil {
				fmt.SendJson(map[string]interface{}{
					"error": "user data invalid"
				})
			}
			// To process the user information
		},
	)
```

## fasthttprouter

[buaazp/fasthttprouter](https://github.com/buaazp/fasthttprouter) forks
[julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) adding
support for the [valyala/fasthttp](https://github.com/valyala/fasthttp).

The implementation is very efficient. However, sometimes we could not find a way
to place our routes the exact way we wanted to. In order to solve this problem,
we implemented our own version (unfortunately less effective).
