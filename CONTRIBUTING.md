## Contributing

### Prerequisites

What things you need to setup the project:

- [go](https://golang.org/doc/install)
- [ginkgo](http://onsi.github.io/ginkgo/)

### Environment

For start developing the SDK you must clone the project:

```bash
git clone git@github.com:lab259/hermes.git
```

Now, the dependencies must be installed.

```bash
go mod download
```

:wink: Finally, you are done to start developing.

### Running tests

In the root directory (where you can find a file named `Makefile`), execute:

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

### Running examples

In the root directory, execute:

```bash
make $EXAMPLE run
```

`$EXAMPLE` is any example listed in `/examples` folder (eg.: `make todos run`)
