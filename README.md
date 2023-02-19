# Helmit

Helmit is a tool for validating, linting, and testing Helm charts. It provides a simple command-line interface for loading charts and running tests.

## Requirements

Helm  must be installed in your environment to use Helmit.
Please make sure to install the appropriate version of Helm before using Helmit. You can find more information about how to install Helm in the [official Helm documentation](https://helm.sh/docs/intro/install/).

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