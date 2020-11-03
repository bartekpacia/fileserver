# Extremely simple fileserver

This tiny program makes it easy to serve any JSON file over the internet (or, more commonly, localhost).

### Serve files

```
+-- project/data
+-- 1
|   +-- data.json
+-- 2
|   +-- data.json
+-- 3
|   +-- and so on...
```

The program will iterate over each JSON file (starting at 1) and serve its contents through TCP/IP. The
content will be refreshed every 5 seconds by default.

### Get help

`go build main.go`

and

`./main --help`

to see command-line options.
