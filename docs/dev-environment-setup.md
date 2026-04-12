# TapX Local Environment Report

Generated: Monday, April 13, 2026 01:16:30

---

## 1. System Information

### OS Details

```
Distributor ID:	Ubuntu
Description:	Ubuntu 24.04.4 LTS
Release:	24.04
Codename:	noble
```

### Kernel

```
Linux banded 6.17.0-20-generic #20~24.04.1-Ubuntu SMP PREEMPT_DYNAMIC Thu Mar 19 01:28:37 UTC 2 x86_64 x86_64 x86_64 GNU/Linux
```

### CPU

```
Architecture:                            x86_64
CPU(s):                                  12
On-line CPU(s) list:                     0-11
Model name:                              Intel(R) Core(TM) i5-10400 CPU @ 2.90GHz
Thread(s) per core:                      2
Core(s) per socket:                      6
CPU(s) scaling MHz:                      19%
NUMA node0 CPU(s):                       0-11
```

### RAM

```
               total        used        free      shared  buff/cache   available
Mem:            15Gi       7.2Gi       5.5Gi       935Mi       4.0Gi       8.2Gi
Swap:          1.9Gi          0B       1.9Gi
```

### Disk Usage

```
Filesystem      Size  Used Avail Use% Mounted on
/dev/nvme0n1p3  142G   21G  115G  16% /
```

### Current User

```
shihab
```

---

## 2. Core Dev Tools

### Go

```
go version go1.26.2 linux/amd64
```

### Docker

```
Docker version 29.4.0, build 9d7ad9f
```

### Docker Compose

```
Docker Compose version v5.1.1
```

### Git

```
git version 2.43.0
```

### VS Code

```
1.115.0
41dd792b5e652393e7787322889ed5fdc58bd75b
x64
```

### curl

```
curl 8.5.0 (x86_64-pc-linux-gnu) libcurl/8.5.0 OpenSSL/3.0.13 zlib/1.3 brotli/1.1.0 zstd/1.5.5 libidn2/2.3.7 libpsl/0.21.2 (+libidn2/2.3.7) libssh/0.10.6/openssl/zlib nghttp2/1.59.0 librtmp/2.3 OpenLDAP/2.6.10
Release-Date: 2023-12-06, security patched: 8.5.0-2ubuntu10.8
Protocols: dict file ftp ftps gopher gophers http https imap imaps ldap ldaps mqtt pop3 pop3s rtmp rtsp scp sftp smb smbs smtp smtps telnet tftp
Features: alt-svc AsynchDNS brotli GSS-API HSTS HTTP2 HTTPS-proxy IDN IPv6 Kerberos Largefile libz NTLM PSL SPNEGO SSL threadsafe TLS-SRP UnixSockets zstd
```

### wget

```
GNU Wget 1.21.4 built on linux-gnu.

-cares +digest -gpgme +https +ipv6 +iri +large-file -metalink +nls
+ntlm +opie +psl +ssl/openssl

Wgetrc:
    /etc/wgetrc (system)
Locale:
    /usr/share/locale
Compile:
    gcc -DHAVE_CONFIG_H -DSYSTEM_WGETRC="/etc/wgetrc"
    -DLOCALEDIR="/usr/share/locale" -I. -I../../src -I../lib
    -I../../lib -Wdate-time -D_FORTIFY_SOURCE=3 -DHAVE_LIBSSL -DNDEBUG
    -g -O2 -fno-omit-frame-pointer -mno-omit-leaf-frame-pointer
    -ffile-prefix-map=/build/wget-LWnKWI/wget-1.21.4=. -flto=auto
    -ffat-lto-objects -fstack-protector-strong -fstack-clash-protection
    -Wformat -Werror=format-security -fcf-protection
    -fdebug-prefix-map=/build/wget-LWnKWI/wget-1.21.4=/usr/src/wget-1.21.4-1ubuntu4.1
    -DNO_SSLv2 -D_FILE_OFFSET_BITS=64 -g -Wall
Link:
    gcc -DHAVE_LIBSSL -DNDEBUG -g -O2 -fno-omit-frame-pointer
    -mno-omit-leaf-frame-pointer
    -ffile-prefix-map=/build/wget-LWnKWI/wget-1.21.4=. -flto=auto
    -ffat-lto-objects -fstack-protector-strong -fstack-clash-protection
    -Wformat -Werror=format-security -fcf-protection
    -fdebug-prefix-map=/build/wget-LWnKWI/wget-1.21.4=/usr/src/wget-1.21.4-1ubuntu4.1
    -DNO_SSLv2 -D_FILE_OFFSET_BITS=64 -g -Wall -Wl,-Bsymbolic-functions
    -flto=auto -ffat-lto-objects -Wl,-z,relro -Wl,-z,now -lpcre2-8
    -luuid -lidn2 -lssl -lcrypto -lz -lpsl ../lib/libgnu.a

Copyright (C) 2015 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later
<http://www.gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Originally written by Hrvoje Niksic <hniksic@xemacs.org>.
Please send bug reports and questions to <bug-wget@gnu.org>.
```

### make

```
GNU Make 4.3
Built for x86_64-pc-linux-gnu
Copyright (C) 1988-2020 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
```

### jq

```
jq-1.7
```

---

## 3. Go Environment

### Go Path Details

```
GOPATH: /home/shihab/go
GOROOT: /snap/go/11127
GOOS: linux
GOARCH: amd64
Binary: /snap/bin/go
```

---

## 4. Docker Status

### Docker Daemon

```
 Server Version: 28.4.0
 Storage Driver: overlay2
 Operating System: Ubuntu Core 24
 Total Memory: 15.45GiB
```

### Running Containers

```
NAMES           IMAGE                STATUS                 PORTS
tapx-app        tapx-lite:latest     Up 3 hours             0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp
tapx-postgres   postgres:15-alpine   Up 3 hours (healthy)   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp
tapx-redis      redis:7              Up 3 hours (healthy)   0.0.0.0:6379->6379/tcp, [::]:6379->6379/tcp
```

### Docker Images

```
REPOSITORY   TAG         SIZE
tapx-lite    latest      19MB
redis        7           117MB
postgres     15-alpine   274MB
```

---

## 5. Network and Ports

### Ports in Use (8080, 5432, 6379, 3000)

```
tcp   LISTEN 0      4096         0.0.0.0:5432       0.0.0.0:*
tcp   LISTEN 0      4096         0.0.0.0:8080       0.0.0.0:*
tcp   LISTEN 0      4096         0.0.0.0:6379       0.0.0.0:*
tcp   LISTEN 0      4096            [::]:5432          [::]:*
tcp   LISTEN 0      4096            [::]:8080          [::]:*
tcp   LISTEN 0      4096            [::]:6379          [::]:*
(empty means ports are free)
```

### Internet Connectivity

```
Internet: CONNECTED
```

---

## 6. Git Configuration

### Git Identity

```
Name: Shihab
Email: shihabud121@outlook.com
```

### SSH Key for GitHub

```
SSH Key found: ~/.ssh/id_ed25519.pub
```

---

## 7. AWS CLI

### AWS CLI Version

```
AWS CLI: NOT INSTALLED (needed later for cloud deployment)
```

---

## 8. Shell and Terminal

### Shell Info

```
Current shell: /bin/zsh
Shell version: zsh 5.9 (x86_64-ubuntu-linux-gnu)
```

### PATH

```
/home/shihab/.npm-global/bin
/usr/local/sbin
/usr/local/bin
/usr/sbin
/usr/bin
/sbin
/bin
/usr/games
/usr/local/games
/snap/bin
/snap/bin
```

---

## 9. Summary

| Tool           | Status        |
| -------------- | ------------- |
| Go             | INSTALLED     |
| Docker         | INSTALLED     |
| Docker Compose | INSTALLED     |
| Git            | INSTALLED     |
| VS Code        | INSTALLED     |
| curl           | INSTALLED     |
| AWS CLI        | NOT INSTALLED |
| jq             | INSTALLED     |
| make           | INSTALLED     |

---
