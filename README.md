# go-diskfs-playground

[kernel/initrd/rootfs](https://github.com/takumin/ubuntu-onmemory-image)

```yaml:meta-data
instance-id: iid-cloudimg
local-hostname: cloudimg
```

```yaml:user-data
#cloud-config

users:
  - name: ubuntu
    gecos: Ubuntu
    shell: /bin/bash
    home: /home/ubuntu
    lock_passwd: False
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: [adm, users, sudo]
    passwd: $6$KKyD4e8i$CYC4kPov1t5hKdexIYHzEaRsIe.ZQYXD8AzCASmyPhPCUpJ5tLzfs99YSAbtP96Y03FaJi0T.vj8zBHrDx0b51
```

```sh
cloud-hypervisor --kernel /usr/local/share/hypervisor-fw/hypervisor-fw --disk path=/tmp/disk.img --rng
```
