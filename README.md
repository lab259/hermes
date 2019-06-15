[![CircleCI](https://circleci.com/gh/lab259/hermes.svg?style=shield)](https://circleci.com/gh/lab259/hermes)
[![codecov](https://codecov.io/gh/lab259/hermes/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/hermes)
[![GoDoc](https://godoc.org/github.com/lab259/hermes?status.svg)](http://godoc.org/github.com/lab259/hermes)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/hermes)](https://goreportcard.com/report/github.com/lab259/hermes)

# HTTP

This is one more web development framework for GO using the fasthttp as base.
Or, maybe, it is just a set of utility of helpful functions that apply some
sugar to the fasthttp "daily" usage.

# Extra Features

- [Grouping](#grouping)
- [Middlewares](#middlewares)
- [Thin JSON Layer](#thin-json-layer)

## Routing

`Routable` is described as an interface which have the Delete, Get, Head,
Options, Patch, Post, Put methods for routing. This aims to [replace the
fasthttprouter implementation](#fasthttprouter).

In order to get things going, the below example shows how to define an endpoint
serving a GET method:

```go
router := hermes.DefaultRouter()

router.Get("/api/v1/user", func(req hermes.Request, res hermes.Response) hermes.Result {
	res.Data("the user is: Snake Eyes")
})
```

### Grouping

When dealing with routes, groups are awesome!

```go
router := hermes.DefaultRouter()

apiGroup := router.Prefix("/api") // This retuns a `Routable` that can be used
                                  // to create other subgroups or define routes.

apiv1 := apiGroup.Prefix("/v1")   // This is what we define as a subgroup.
apiv1.Get(                        // Now a definition of the route itself.
	"/user",
	func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data("the user is: Snake Eyes")
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
router := hermes.DefaultRouter()

apiGroup := router.Prefix("/api")

apiv1 := apiGroup.Prefix("/v1")

apiv1.Use(func(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
	// This is a middleware that could do something smart...
	return next(req, res)
})

apiv1.Get(
	"/user",
	func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data("the user is: Snake Eyes")
	},
)
```

Yet, you can also define multiple middlewares for each route and their priority
will be from the left to the right.

```go
router := hermes.DefaultRouter()

apiGroup := router.Prefix("/api")
apiv1 := apiGroup.Prefix("/v1")
apiv1.With(
	func(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
		// This is a middleware that could do something smart...
		return next(ctx)
	},
	func(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
		// This is a second middleware for the endpoint...
		return next(ctx)
	},
).Get(
	"/user",
	func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data("the user is: Snake Eyes")
	},
)
```

Middlewares are also supported on groups:

```go
router := hermes.DefaultRouter()

apiGroup := router.Prefix("/api").With(func(req hermes.Request, res hermes.Resonse, next hermes.Handler) hermes.Result {
	// This is a middleware that could do something smart...
	return next(ctx)
})
apiv1 := apiGroup.Prefix("/v1").With(func(req hermes.Request, res hermes.Resonse, next hermes.Handler) hermes.Result {
	// Yet another middleware applied just for this subgroup...
	return next(ctx)
})
apiv1.With(func(req hermes.Request, res hermes.Resonse, next hermes.Handler) hermes.Result {
	// This is a middleware that is applied just for this endpoint
	return next(ctx)
}).Get(
	"/user",
	func(req hermes.Request, res hermes.Resonse) hermes.Result {
		return res.Data("the user is: Snake Eyes")
	},
)
```

Again, there is no performance difference when using middlewares in a specific
route or in a whole group. The internal implementation will append both
middleware definitions into one big sequence of middlewares for each route.

### Thin JSON layer

For simple sake of ease the use of sending and receiving JSON objects `res.Data`
and `req.Data` methods were added.

#### Sending a JSON

The following is an example of sending a JSON document:

```go
router := hermes.DefaultRouter()

apiv1 := router.Prefix("/api/v1")
apiv1.Get(
	"/user",
	func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data(map[string]interface{}{
			"name": "Snake Eyes",
			"email": "s.eyes@gijoe.com",
		})
	},
)
```

### Receiving a JSON

The following is an example of receiving a JSON document:

```go
router := hermes.DefaultRouter()

apiv1 := router.Group("/api/v1")
apiv1.Post(
	"/user",
	func(req hermes.Request, res hermes.Response) hermes.Result {
		user := make(map[string]interface{})
		if err := req.Data(&user); err != nil {
			return res.Status(400).Data(map[string]interface{}{
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
