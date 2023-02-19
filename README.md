# Helmit

Helmit is a tool for validating, linting, and testing Helm charts. It provides a simple command-line interface for loading charts and running tests.

## Usage

To validate and lint a chart:

`$ helmit --chart path/to/chart`

To initialize a Kubernetes test environment:

`$ helmit --inittestenv path/to/kubeconfig`

To test a chart, run:

`$ helmit --chart path/to/chart --test`

This will deploy the chart to a Kubernetes environment configured using the `--inittestenv` flag.


## Flags

- `-c, --chart string`: Path to the Helm chart
- `-h, --help`: Output usage information
- `--inittestenv string`: Initialize a test environment with the given kubeconfig file
- `-t, --test`: Test the Helm chart after loading


## License

This project is licensed under the [MIT License](LICENSE).