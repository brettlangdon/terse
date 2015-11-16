terse
=====
Short URL server written in [go](https://golang.org) which stores all URL short code mappings in an in-memory least recently used (LRU) cache.

## Getting Started
```bash
# Install the `terse` command
go get github.com/brettlangdon/terse/cmd/...

# Start the server
terse
```

```
$ curl -i -X POST -d "https://github.com/brettlangdon/terse" http://127.0.0.1:5893
HTTP/1.1 201 Created
Date: Mon, 16 Nov 2015 01:10:12 GMT
Content-Length: 28
Content-Type: text/plain; charset=utf-8

http://127.0.0.1:5892/DEtr1b
$ curl -i http://127.0.0.1:5892/DEtr1b
HTTP/1.1 301 Moved Permanently
Location: https://github.com/brettlangdon/terse
Date: Mon, 16 Nov 2015 01:10:54 GMT
Content-Length: 72
Content-Type: text/html; charset=utf-8

<a href="https://github.com/brettlangdon/terse">Moved Permanently</a>.

```

## Installing
Install via `go get` with `go get github.com/brettlangdon/terse/cmd/...`.

`terse` requires `GO15VENDOREXPERIMENT=1` in order to build.

## Usage
```
$ ./terse --help
usage: terse [--max MAX] [--bind BIND] [--server SERVER]

options:
  --max MAX, -m MAX      max number of links to keep ("0" means no limit) [default: 1000]
  --bind BIND, -b BIND   "[host]:<port>" to bind the server to [default: 127.0.0.1:5892]
  --server SERVER, -s SERVER
                         base server url to generate links as (e.g. "https://short.domain.com") [default: "http://<bind>"]
```

### Max
`terse` uses an in-memory least recently used (LRU) cache and defaults to having a limit of only 1000 links. This means that after generating 1000 URLs the least recently used URLs will start to be purged from the cache.

You can control this limit with the `--max` parameter. Setting the value to `0` will disable the size limit of the cache.

### Bind
If you would like to change the host or port that `terse` listens on by default then supply the `--bind` parameter. Examples: `127.0.0.1:5892`, `:80`, `0.0.0.0:8000`, etc.

By default `terse` will listen to `127.0.0.1:5892`.

### Server
By default `terse` will respond with URLs generated based on the `--bind` parameter. This means that by default you will get responses like `http://127.0.0.1:5892/<SHORT CODE>`. Instead, if you would like to generate URLs with a different scheme and hostname, supply the correct `--server` parameter.

For example:

```bash
terse --server "https://s.example.org"
```

Will produce URLs like `https://s.example.org/<SHORT CODE>`.

## API
`terse` has a very simple to use API.

### Creating Links
To create a new short link, issue a `POST` request to the server with the `POST` body of the URL to shorten.

For example:

```bash
curl -X POST -d "https://github.com/brettlangdon/terse" http://127.0.0.1:5892
```

A valid response will have the status `201 Created`, and the response body will be the short URL.

You may also receive a `400 Bad Request` if the `POST` body is not a valid URL. Or a `409 Conflict` if there was a hash conflict (two different URLs ended up with the same short code).

### Accessing Links
To utilize a short link, issues a `GET` request for the URL. If the short code exists in the system then the response will be a `301 Moved Permanently` for the stored URL.

If the short code does not exist, then you will receive a `404 Not Found` response.
