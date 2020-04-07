# Trieres

[![Build Status](https://cloud.drone.io/api/badges/jakolehm/trieres/status.svg)](https://cloud.drone.io/jakolehm/trieres)

Trieres is a [K3s](https://github.com/rancher/k3s) cluster lifecycle management tool.

## Usage


### Install/re-configure/upgrade cluster

```
$ trieres up --config ./cluster.yml
```

### Fetch admin kubeconfig from cluster

```
$ trieres kubeconfig --config ./cluster.yml
```

## Example cluster.ymls

### Minimal cluster.yml example

```yaml
hosts:
  - address: 1.2.3.4
    role: master
```

### Full cluster.yml example

```yaml
token: verysecret
version: "v1.17.2+k3s1"
hosts:
  - address: "1.2.3.4"
    role: master
    user: root
    sshKeyPath: "~/.ssh/id_rsa"
    sshPort: 22
    extraArgs: []
  - address: "2.3.4.5"
    role: worker
    user: root
    sshKeyPath: "~/.ssh/id_rsa"
    sshPort: 22
    extraArgs:
      - "--node-label foo=bar"
manifests:
  - ./manifests/*.yaml
```
