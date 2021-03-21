# Payment Gateway Exercise - Payment Processor Service

Please check out documentation for the whole system [here](https://github.com/gustavooferreira/pgw-docs).

This repository structure follows [this convention](https://github.com/golang-standards/project-layout).

---

## Tip

> If you run `make` without any targets, it will display all options available on the makefile followed by a short description.

# Build

To build a binary, run:

```bash
make build
```

The `api-server` binary will be placed inside the `bin/` folder.

---

# Tests

To run tests:

```bash
make test
```

To get coverage:

```bash
make coverage
```

# Docker

To build the docker image, run:

```bash
make build-docker
```

The docker image is named `pgw/payment-processor-api-server`.

Create a file with some credit cards and reasons to fail, like this:

```yaml
credit_cards:
  4000000000000119: "authorisation failure"
  4000000000000259: "capture failure"
  4000000000003238: "refund failure"
```

And start a docker container like this:

```bash
docker run --rm --name pgw-payment-processor-service -p 127.0.0.1:9000:8080/tcp -v "$(pwd)"/credit_cards.yaml:/credit_cards.yaml:ro -e PGW_PAYMENT_PROCESSOR_APP_OPTIONS_CREDITCARDS_FILENAME=/credit_cards.yaml pgw/payment-processor-api-server
```

This assumes the yaml file created is called `credit_cards.yaml` and is placed in the current directory.

Once the container is running, you can make a request like this:

```bash
curl -i -X POST http://localhost:9000/api/v1/authorise -d '{"":""}'
```

# Design

