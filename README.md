# strmatch
search for regexp in stdin

#### usage example

```bash
$ ls -l
total 41
-rw-r--r-- 1 xxxxx 1049089 35147 Jul  19 17:35 LICENSE
-rw-r--r-- 1 xxxxx 1049089  1755 Jul  19 17:42 main.go
-rw-r--r-- 1 xxxxx 1049089    38 Jul  19 17:35 README.md

$ ls -l | strmatch  "\d\d ([^ ]+)$"
LICENSE
main.go
README.md

$ ls -l | strmatch -sep " " "\d\d ([^ ]+)$"
LICENSE main.go README.md
```
