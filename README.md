# vet-k8sdepr
Code checker for Kubernetes API deprecation/removal.


# Install

```
go get -u github.com/yoichiwo7/vet-k8sdepr
```

# Usage

```
vet-k8sdepr -targetVersion VERSION [-flag] [package]

Flags:
  -targetVersion string
        target semantic version of the Kubernetes (ex. v1.16.0)
  -ignoreDeprecation
        ignore deprecation detection
  -ignoreRemoval
        ignore removal detection
```

# Example

```bash
# Check codes in current directory against Kubernetes v1.17.0
vet-k8sdepr -targetVersion v1.17.0 ./...
```