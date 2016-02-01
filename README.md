# goenvtemplator
Tool to template configuration files by environment variables and optionally replace itself with the target binary.

goenvtemplator is a simple app, that can template your config files by environment variables and optionally replace itself (by exec syscall) with the application binary. So at the end your application runs directly under the process that run this tool like docker as if it was originally the entrypoint itself.

This tool is ideal for use without polluting you environment with dependencies. It is fully statically linked so it has no dependencies whatsoever. If you use Dockerfile you don't even need wget or curl since it can be installed only by dockerfile's ADD instruction. 

## Installation
wget
```bash
wget https://github.com/seznam/goenvtemplator/releases/download/v1.0.0/goenvtemplator-amd64 -O /usr/local/bin/goenvtemplator
chmod +x /usr/local/bin/goenvtemplator
```

Dockerfile
```Dockerfile
ADD https://github.com/seznam/goenvtemplator/releases/download/v1.0.0/goenvtemplator-amd64 /usr/local/bin/goenvtemplator
RUN chmod +x /usr/local/bin/goenvtemplator
```


## Usage
goenvtemplator -help
```
Usage of goenvtemplator:
  -debug-templates
    	Print processed templates to stdout.
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
goenvtemplator -template /path/to/server.conf.tmpl:/path/to/server.conf  -template /path/to/server2.conf.tmpl:/path/to/server2.conf
```

### Dockerfile
```Dockerfile
ENTRYPOINT ["/usr/local/bin/goenvtemplator", "-template", "/path/to/server.conf.tmpl:/path/to/server.conf", "-exec"]
CMD ["/usr/bin/server-binary", "server-argument1", "server-argument2", "..."]
```

## Using Templates
Templates use Golang [text/template](http://golang.org/pkg/text/template/).

### Built-in functions
There are a few built in functions as well:
  * `env "ENV_NAME"` - Accesses environment variables. If it does not exist return empty string. `{{ env "TIMEOUT_MS }}`
  * `require (env "ENV_NAME")` - Renders an error if environments variable does not exists. If it is equal to empty string, returns empty string.  `{{ require (env "TIMEOUT_MS) }}`
  * `default $default_1 $default_2 $default_3` - Returns a first argument that exists. If none is valid it generates error `{{ default (env "SPECIFIC_TIMEOUT_MS") (env "GENERAL_TIMEOUT_MS") "1000" }}`
