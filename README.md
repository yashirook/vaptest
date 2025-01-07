# vaptest

**vaptest** is a minimal command-line tool for testing Kubernetes `ValidationAdmissionPolicy` and `ValidationAdmissionPolicyBinding` against actual Kubernetes manifests.

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
$ make build

# Check Version
$ ./bin/vaptest version
vaptest version: 0.1.0 (commit ad669bc, build at 2025-01-07T22:43:28Z)
```

Please move the binary file to a directory that is included in the system PATH.
```bash
$ mv ./bin/vaptest /usr/local/bin
```

## Usage
### Validate Manifests
Validate a single Kubernetes manifest against your policies:

```bash
$ vaptest validate --policies=./example/policy/policy.yaml --targets=./example/target/valid-deployment.yaml
all validation success!
```

Validate all manifests in a directory:

```bash
$ vaptest validate --policies=./example/policy/policy.yaml --targets=./example/target/valid-deployment.yaml
POLICY         EVALUATED_RESOURCE            RESULT  ERRORS
require-label  deployments/nginx-deployment  Fail    Deployment has to have namespace (Expression: has(object.metadata.namespace))
require-label  services/nginx-service        Fail    Deployment has to have label (Expression: has(object.metadata.labels))
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