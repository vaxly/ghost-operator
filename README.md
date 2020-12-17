# ghost-operator

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
      secretName: blog-vaxly-io-prod
    hosts:
      - blog.vaxly.io
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-prod
      kubernetes.io/tls-acme: "true"
  config:
    url: https://blog.vaxly.io
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
        password: secret
        database: ghostdb
```

Note enabling tls assumes you gave cert-manager already installed and configured

