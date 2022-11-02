# Snowputer
> **Fast Golang Tool take a filename with urls and test if the put method is enabled it upload a file **


# Install
```
$ go install github.com/yghonem14/snowputer@latest
```

## Basic Usage
Snowputer accepts only text files with -f option:

```
$ snowputer -f ~/targets/urls.txt
owkeoijwrijw.server.com -> vulnerable
```


## Concurrency

You can set the concurrency value with the `-c` flag:

```
$ snowputer -f ~/targets/urls.txt -c 35
```
