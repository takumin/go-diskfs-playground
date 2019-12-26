# go-diskfs-playground

[kernel/initrd/rootfs](https://github.com/takumin/ubuntu-onmemory-image)

```yaml:meta-data
instance-id: iid-cloudimg
local-hostname: cloudimg
```

```yaml:user-data
#cloud-config
password: ubuntu
chpasswd: { expire: False }
ssh_pwauth: True
```
