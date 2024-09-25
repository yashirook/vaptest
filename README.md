# vaptest

**vaptest** is a command-line tool for testing Kubernetes `ValidationAdmissionPolicy` and `ValidationAdmissionPolicyBinding` against actual Kubernetes manifests.

> **Note:** This project is currently under development. Features and interfaces may change.

## What is ValidationAdmissionPolicy?
[ValidationAdmissionPolicy](https://kubernetes.io/docs/reference/access-authn-authz/validating-admission-policy/) is a Kubernetes feature that allows cluster administrators to define custom admission control policies using the Common Expression Language (CEL). It enables the enforcement of complex validation logic on Kubernetes resources during the admission phase, without the need to write custom webhooks.

## Features

- Validate Kubernetes manifests against defined `ValidationAdmissionPolicy` rules.
- Simulate admission evaluations locally without deploying to a cluster.
- Designed for integration in local development and CI pipelines.
- Output detailed validation results in text or JSON format.

## Installation

To install `vaptest`, make sure you have [Go](https://golang.org/dl/) installed (version 1.16 or higher), and then run:

```bash
go install github.com/yashirook/vaptest@latest
```

## Usage
### Validate Manifests
Validate a single Kubernetes manifest against your policies:

```bash
vaptest validate --policy-dir=./policies --manifest=./manifests/deployment.yaml
```

Validate all manifests in a directory:

```bash
vaptest validate -p ./policies -m ./manifests/
```

### Run Test Cases
Execute test cases defined in a directory:

```bash
vaptest test --policy-dir=./policies --test-dir=./tests/
```

## Command-Line Options
- --policy-dir, -p: Directory containing ValidationAdmissionPolicy and ValidationAdmissionPolicyBinding definitions (required).
- --manifest, -m: Kubernetes manifest file or directory to validate.
- --test-dir, -t: Directory containing test cases.
- --output, -o: Output format (text or json). Default is text.
- --verbose, -v: Enable verbose output.

## Examples
Validate a manifest with verbose output:

```bash
vaptest validate -p ./policies -m ./manifests/deployment.yaml -v
```

Output results in JSON format:

```bash
vaptest validate -p ./policies -m ./manifests/ -o json
```

## Development Status
This project is in active development. Some features may not be fully implemented, and the interface is subject to change. Contributions and feedback are welcome!

## Contributing
Fork the repository.
Create a new feature branch (git checkout -b feature/my-feature).
Commit your changes (git commit -am 'Add new feature').
Push to the branch (git push origin feature/my-feature).
Open a Pull Request.

## License
This project is licensed under the Apache License 2.0. See the LICENSE file for details.

## Contact
For questions or suggestions, please open an issue or contact yashiro.kentaro@gmail.com.