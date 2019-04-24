## Contributing

### Prerequisites

What things you need to setup the project:

- [go](https://golang.org/doc/install)
- [golang/dep](https://github.com/golang/dep)
- [ginkgo](http://onsi.github.io/ginkgo/)

### Environment

For start developing the SDK you must create a fake `GOPATH` structure:

```
+-- /
|---- src
|------ github.com
|-------- lab259
|---------- http <- Here is where you will clone this repository.
```

Use the following command:

```bash
mkdir -p src/github.com/lab259/http && git clone git@github.com:lab259/http.git src/github.com/lab259/http
```

Now, the dependencies must be installed.

```
cd src/github.com/lab259/http && make dep-ensure
```

:wink: Finally, you are done to start developing.

### Running tests

In the `src/github.com/lab259/http` directory, execute:

```bash
make test
```

To enable coverage, execute:

```bash
make coverage
```

To generate the HTML coverage report, execute:

```bash
make coverage coverage-html
```
