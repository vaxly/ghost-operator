# ghost-operator

This operator deploy both ghost and mysql with auto configuration.

## Installation

First apply the manifest
```shell
kubectl apply -n vaxly -f https://raw.githubusercontent.com/vaxly/ghost-operator/main/bundle/manifests/ghost.vaxly.io_blogs.yaml
```

Then copy the file below, append it to your needs, then apply it to your cluster
Example values 
```yaml
apiVersion: ghost.vaxly.io/v1alpha1
kind: Blog
metadata:
  name: vaxly-blog
spec:
  replicas: 1
  image: ghost:3.40.1
  persistent:
    enabled: true
    size: 10Gi
  serviceType: ClusterIP
  ingress:
    enabled: true
    tls:
      enabled: true
      secretName: blog-vaxly-io-prod  # Must change
    hosts:
      - blog.vaxly.io # Must change
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-prod # Must change
      kubernetes.io/tls-acme: "true"
  config:
    url: https://blog.vaxly.io  # Must change
    server:
      host: 0.0.0.0
      port: 80
    mail:
      transport: direct
    logging:
      transports:
        - stdout
    database:
      client: mysql
      connection:
        host: localhost
        port: 3306
        user: root
        password: secret  # Must change
        database: ghostdb
```

Note enabling tls assumes you gave cert-manager already installed and configured

