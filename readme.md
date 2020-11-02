# auto ssh

## cmd
- `assh`

## config
- `~/.ssh/assh.config`
- `format: yaml`
```
- name: test f28
  host: 192.168.112.2
  user: root
  port: 22
  password: pass
  referer: ""
  cmd: ""
- name: test1
  host: 192.168.112.3
  user: root
  port: 22
  password: pass
  referer: ""
  cmd: ""
```