# Additional Information

## Module dependencies

If you’re creating a package or application which can be downloaded and used by other
people and programs, then it’s good practice for your module path to equal the location that the code can be downloaded from.
For instance, if your package is hosted at `https://github.com/foo/bar` then the module
path for the project should be github.com/foo/bar.

## Web Application Basics

### Network Addresses

The TCP network address that you pass to `http.ListenAndServe()` should be in the format
"host:port" . If you omit the host (like we did with ":4000" ) then the server will listen on all
your computer’s available network interfaces. Generally, you only need to specify a host in
the address if your computer has multiple network interfaces and you want to listen on just
one of them.
In other Go projects or documentation you might sometimes see network addresses written
using named ports like ":http" or ":http-alt" instead of a number. If you use a named
port then Go will attempt to look up the relevant port number from your `/etc/services` file
when starting the server, or will return an error if a match can’t be found.

### Using go run

During development the go run command is a convenient way to try out your code. It’s
essentially a shortcut that compiles your code, creates an executable binary in your /tmp
directory, and then runs this binary in one step.
It accepts either a space-separated list of .go files, the path to a specific package (where
the . character represents your current directory), or the full module path. For our
application at the moment, the three following commands are all equivalent:

```bash
$ go run main.go
$ go run .
$ go run 'module/name'
```

## Routing Requests

### Servemux Features and Quirks

- In Go’s servemux, longer URL patterns always take precedence over shorter ones. So, if a
  servemux contains multiple patterns which match a request, it will always dispatch the
  request to the handler corresponding to the longest pattern. This has the nice side-
  effect that you can register patterns in any order and it won’t change how the servemux
  behaves.
- Request URL paths are automatically sanitized. If the request path contains any . or ..
  elements or repeated slashes, the user will automatically be redirected to an equivalent
  clean URL. For example, if a user makes a request to /foo/bar/..//baz they will
  automatically be sent a 301 Permanent Redirect to /foo/baz instead.

- If a subtree path has been registered and a request is received for that subtree path
  without a trailing slash, then the user will automatically be sent a
  301 Permanent Redirect to the subtree path with the slash added. For example, if you
  have registered the subtree path /foo/ , then any request to /foo will be redirected to
  /foo/.

### Host Name Matching

It’s possible to include host names in your URL patterns. This can be useful when you want
to redirect all HTTP requests to a canonical URL, or if your application is acting as the back
end for multiple sites or services. For example:

```go
mux := http.NewServeMux()
mux.HandleFunc("foo.example.org/", fooHandler)
mux.HandleFunc("bar.example.org/", barHandler)
mux.HandleFunc("/baz", bazHandler)
```

When it comes to pattern matching, any host-specific patterns will be checked first and if
there is a match the request will be dispatched to the corresponding handler. Only when
there isn’t a host-specific match found will the non-host specific patterns also be checked.

### RESTful Routing?

It’s important to acknowledge that the routing functionality provided by Go’s servemux is
pretty lightweight. It doesn’t support routing based on the request method, it doesn’t
support semantic URLs with variables in them, and it doesn’t support regexp-based
patterns. If you have a background in using frameworks like Rails, Django or Laravel you
might find this a bit restrictive… and surprising!
But don’t let that put you off. The reality is that Go’s servemux can still get you quite far,
and for many applications is perfectly sufficient. For the times that you need more, there’s a
huge choice of third-party routers that you can use instead of Go’s servemux. We’ll look at
some of the popular options later in the book.

### Manipulating the Header Map

In the code above we used w.Header().Set() to add a new header to the response header
map. But there’s also Add() , Del() and Get() methods that you can use to read and
manipulate the header map too.

```go
// Set a new cache-control header. If an existing "Cache-Control" header exists
// it will be overwritten.
w.Header().Set("Cache-Control", "public, max-age=31536000")
// In contrast, the Add() method appends a new "Cache-Control" header and can
// be called multiple times.
w.Header().Add("Cache-Control", "public")
w.Header().Add("Cache-Control", "max-age=31536000")
// Delete all values for the "Cache-Control" header.
w.Header().Del("Cache-Control")
// Retrieve the first value for the "Cache-Control" header.
w.Header().Get("Cache-Control")
```

### System-Generated Headers and Content Sniffing

When sending a response Go will automatically set three system-generated headers for you:
Date and Content-Length and Content-Type .
The Content-Type header is particularly interesting. Go will attempt to set the correct one
for you by content sniffing the response body with the http.DetectContentType() function.
If this function can’t guess the content type, Go will fall back to setting the header
Content-Type: application/octet-stream instead.
The http.DetectContentType() function generally works quite well, but a common gotcha
for web developers new to Go is that it can’t distinguish JSON from plain text. So, by
default, JSON responses will be sent with a Content-Type: text/plain; charset=utf-8
header. You can prevent this from happening by setting the correct header manually like
so:

```go
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"name":"Alex"}`))
```

### Header Canonicalization

When you’re using the Add() , Get() , Set() and Del() methods on the header map, the
header name will always be canonicalized using the textproto.CanonicalMIMEHeaderKey()
function. This converts the first letter and any letter following a hyphen to upper case, and
the rest of the letters to lowercase. This has the practical implication that when calling
these methods the header name is case-insensitive.
If you need to avoid this canonicalization behavior you can edit the underlying header map
directly (it has the type map[string][]string ). For example:

`w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}`

- Note: If a HTTP/2 connection is being used, Go will always automatically convert the
  header names and values to lowercase for you as per the HTTP/2 specifications.

### Suppressing System-Generated Headers

The Del() method doesn’t remove system-generated headers. To suppress these, you need
to access the underlying header map directly and set the value to nil . If you want to
suppress the Date header, for example, you need to write:

`w.Header()["Date"] = nil`

## URL Query Strings