# Patches

## cilikube-web-monitoring-rbac-storage.patch

Cloud agent token has **read-only** access to `cillianxtech/cilikube-web` (push returns 403).

Apply the Vue frontend updates locally:

```bash
git clone https://github.com/cillianxtech/cilikube-web.git
cd cilikube-web
git checkout -b cursor/sync-monitoring-rbac-storage-b44a
git apply ../cilikube/patches/cilikube-web-monitoring-rbac-storage.patch
# or: git am < ../cilikube/patches/cilikube-web-monitoring-rbac-storage.patch
git push -u origin HEAD
```

Contents: monitoring page, StorageClass, K8s RBAC views, clusterId injection, locales.
