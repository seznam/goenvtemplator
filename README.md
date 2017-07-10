# goenvtemplator
Tool to template configuration files by environment variables and optionally replace itself with the target binary.

goenvtemplator is a simple app, that can template your config files by environment variables and optionally replace itself (by exec syscall) with the application binary. So at the end your application runs directly under the process that run this tool like docker as if it was originally the entrypoint itself.

This tool is ideal for use without polluting you environment with dependencies. It is fully statically linked so it has no dependencies whatsoever. If you use Dockerfile you don't even need wget or curl since it can be installed only by dockerfile's ADD instruction. 

## Installation
wget
```bash
wget https://github.com/seznam/goenvtemplator/releases/download/v2.0.0-rc3/goenvtemplator-amd64 -O /usr/local/bin/goenvtemplator2
chmod +x /usr/local/bin/goenvtemplator2
```

Dockerfile
```Dockerfile
ADD https://github.com/seznam/goenvtemplator/releases/download/v2.0.0-rc3/goenvtemplator-amd64 /usr/local/bin/goenvtemplator2
RUN chmod +x /usr/local/bin/goenvtemplator2
```

## Building from source
```bash
# if you have glide already get the binary
go get github.com/Masterminds/glide
# install dependencies
$GOPATH/bin/glide i
make
```

## Usage
goenvtemplator2 -help
```
Usage of goenvtemplator:
  -debug-templates
    	Print processed templates to stdout.
  -env-file value
        Additional file with environment variables. Can be passed multiple times. (default [])
  -exec
    	Activates exec by command. First non-flag arguments is the command, the rest are it's arguments.
  -template value
    	Template (/template:/dest). Can be passed multiple times. (default [])
  -v int
    	Verbosity level.
  -version
    	Prints version.
```

### Example
```bash
goenvtemplator2 -template /path/to/server.conf.tmpl:/path/to/server.conf  -template /path/to/server2.conf.tmpl:/path/to/server2.conf
```

### Dockerfile
```Dockerfile
ENTRYPOINT ["/usr/local/bin/goenvtemplator2", "-template", "/path/to/server.conf.tmpl:/path/to/server.conf", "-exec"]
CMD ["/usr/bin/server-binary", "server-argument1", "server-argument2", "..."]
```

### Env files
It is possible to add additional environment variables in multiple env-files.
Existing variables are **not** overwritten.
Environment variables in files can be formated using shell syntax or yaml syntax.

Let us consider an environment file `myenvfile` bellow.
```bash
# cat myenvfile
A=a
B=b
B=bb
#B=bbb
export C=c
# yaml syntax
D: d
```

The behaviour of env-file argument and env variables evaluation is as follows:
```
goenvtemplator2 -env-file myenvfile -exec sh -c 'echo $A'
> a
export A=foo
goenvtemplator2 -env-file myenvfile -exec sh -c 'echo $A'
> foo
goenvtemplator2 -env-file myenvfile -exec sh -c 'echo $B'
> bb
goenvtemplator2 -env-file myenvfile -exec sh -c 'echo $C'
> c
goenvtemplator2 -env-file myenvfile -exec sh -c 'echo $D'
> d
```

## Using Templates
Templates use Golang [text/template](http://golang.org/pkg/text/template/)
and [Sprig](https://github.com/Masterminds/sprig) library.

### Built-in functions
There are a few built in functions as well:
  * `require (env "ENV_NAME")` - Renders an error if environments variable does not exists. If it is equal to empty string, returns empty string.  `{{ require (env "TIMEOUT_MS) }}`
