Go code analyzer for Kubernetes API deprecation/removal.

As Kubernetes version goes up, the resource API such as `apps/v1beta2:Deployment` will be deprecated/removed.
`k8sdepr` detects these API deprecation/removal in Go codes; based on the specified target Kubernetes version.

If you are interested in detecting API deprecation/removal in YAML manifests and Kubernetes cluster resources, then check follwing tools:

- [FairwindsOps/pluto](https://github.com/FairwindsOps/pluto)

# Install

```
go get -u github.com/yoichiwo7/k8sdepr
```

# Usage

- `-targetVersion` falg must be set. The value must follow semantic version. (ex. `v1.16.0`)

```
k8sdepr -targetVersion VERSION [-flag] [package]

Flags:
  -targetVersion string
        target semantic version of the Kubernetes (ex. v1.16.0)
  -ignoreDeprecation
        ignore deprecation detection
  -ignoreRemoval
        ignore removal detection
```

# Example

Run the analyzer command on target source directory. As described in previous section, `-targetVersion` must be set with valid semantic version.

```bash
# Check codes in current directory against Kubernetes v1.17.0
k8sdepr -targetVersion v1.17.0 ./...
```

If the analyzer detects API deprecation or removal, it will prints message like following.

```
/tmp/src/services.go:9:10: apps/v1beta2:Deployment is removed in v1.16.0. Migrate to apps/v1:Deployment.
/tmp/src/ingress.go:38:10: extensions/v1beta1:Ingress is deprecated in v1.14.0. Migrate to networking.k8s.io/v1beta1:Ingress.
```

## Check Non-Module Codes

If you want to check non-module codes, you must setup go module and sync vendor directory first.
The following example shows how to check these code by using old version of `spotahome/redis-operator` which has some deprecated/removed APIs.

```bash
# Get redis-operator (specify version that has some removed APIs)
git clone -b 0.5.0 --depth=1 https://github.com/spotahome/redis-operator
cd redis-operator

# Setup go module
go mod init github.com/spotahome/redis-operator

# Sync vendor directory
go mod vendor

# Check against v1.15.0 -> Detects lots of API deprecation
k8sdepr -targetVersion v1.15.0 ./...

# Check against v1.16.0 -> Detects lots of API removal
k8sdepr -targetVersion v1.16.0 ./...
```