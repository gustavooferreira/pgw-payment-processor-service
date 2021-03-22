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
creditCards:
  4000000000000119: "authorise fail"
  4000000000000259: "capture fail"
  4000000000003238: "refund fail"
```

And start a docker container like this:

```bash
docker run --rm --name pgw-payment-processor-service -p 127.0.0.1:9000:8080/tcp -v "$(pwd)"/edge_cases_credit_cards.yaml:/edge_cases_credit_cards.yaml:ro -e PGW_PAYMENT_PROCESSOR_APP_OPTIONS_CREDITCARDS_FILENAME=/edge_cases_credit_cards.yaml pgw/payment-processor-api-server
```

This assumes the yaml file created is called `edge_cases_credit_cards.yaml` and is placed in the current directory.

Once the container is running, you can make a request like this:

```bash
curl -i -X POST http://localhost:9000/api/v1/authorise -d '{"credit_card": {"name":"customer1", "number": 4000000000000118, "expiry_month":10, "expiry_year":2030, "cvv":123}, "currency": "EUR", "amount": 10.50}'
```

# Design

This service serves as a light dependency for the payment gateway service. Its purpose is to mimic the behavior of a potential payment processor.

This service is pretty dumb, it doesn't check for any constraints and it doesn't rely on a database. In fact, this service will reply successfully to all calls made to it except for a specific number of credit cards specified in a yaml file.

This service reads credit cards and their reason to fail from a yaml file.

This service holds current state in memory which means once restarted all that state is lost, and that's fine for testing purposes.

This service is also responsible for generating a `UID` for each `authorisation` call, as I'm assuming that's how it works in the real world.

It's not meant to be production ready by any means. I'm not using a database to simplify the service, as this service was only created so the payment gateway can simulate talking to an external system to process the payment.

The OpenAPI spec is located in the `openapi` folder.

To view the spec in the Swagger UI [click this link](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/gustavooferreira/pgw-payment-processor-service/master/openapi/spec.yaml).

The requests to this service should go through an authentication/authorization process as well. I have not implemented this to keep the service simple.
