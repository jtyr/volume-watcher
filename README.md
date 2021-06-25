volume-watcher
==============

This is a Kubernetes volume watcher usable as a sidecar. The volume can be
either a `Secret` or a `ConfigMap`. Once the watcher identifies a change in the
volume, it sends an HTTP request to the defined endpoint.


Usage
-----

The sidecar requires the following two environment variables to be set:

- `VOLUMEWATCHER_DIR` - path to the directory where the volume is mounted
- `VOLUMEWATCHER_ENDPOINT` - URL of an endpoint to call when the volume changes

The following `Deployment` example shows how to use it as a sidecar:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agent
spec:
  template:
    spec:
      containers:
        # Sidecar
        - name: volume-watcher
          image: jtyr/volume-watcher
          env:
            - name: VOLUMEWATCHER_DIR
              value: /etc/agent
            - name: VOLUMEWATCHER_ENDPOINT
              value: http://localhost:8080/-/reload
          volumeMounts:
            - mountPath: /etc/agent
              name: agent-config
        # Container which will be reloaded from teh sidecar
        - name: agent
          ports:
            - containerPort: 8080
              name: http-metrics
          ...
          volumeMounts:
            - mountPath: /etc/agent
              name: agent-config
      volumes:
        - configMap:
            name: agent-config
          name: agent-config
```


License
-------

MIT


Author
------

Jiri Tyr
