# goenvtemplator
![GitHub release (latest by date)](https://img.shields.io/github/v/release/seznam/goenvtemplator) ![GitHub all releases](https://img.shields.io/github/downloads/seznam/goenvtemplator/total)


Tool to template configuration files by environment variables and optionally replace itself with the target binary.

goenvtemplator is a simple app, that can template your config files by environment variables and optionally replace itself (by exec syscall) with the application binary. So at the end your application runs directly under the process that run this tool like docker as if it was originally the entrypoint itself.

This tool is ideal for use without polluting you environment with dependencies. It is fully statically linked so it has no dependencies whatsoever. If you use Dockerfile you don't even need wget or curl since it can be installed only by dockerfile's ADD instruction.

## Installation
wget
```bash
wget https://github.com/seznam/goenvtemplator/releases/download/v2.0.0/goenvtemplator2-amd64 -O /usr/local/bin/goenvtemplator2
chmod +x /usr/local/bin/goenvtemplator2
```

Dockerfile
```Dockerfile
ADD https://github.com/seznam/goenvtemplator/releases/download/v2.0.0/goenvtemplator2-amd64 /usr/local/bin/goenvtemplator2
RUN chmod +x /usr/local/bin/goenvtemplator2
```

## Building from source
```bash
make
chmod +x goenvtemplator2
cp goenvtemplator2 /usr/local/bin
```

## Usage
goenvtemplator2 -help
```
Usage of goenvtemplator2:
  -debug-templates
        Print processed templates to stdout.
  -delim-left string
        Override default left delimiter {{.
  -delim-right string
        Override default right delimiter }}.
  -env-file value
        Additional file with environment variables. Can be passed multiple times.
  -exec
        Activates exec by command. First non-flag arguments is the command, the rest are it's arguments.
  -template value
        Template (/template:/dest). Can be passed multiple times.
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
E=$A
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
goenvtemplator2 -env-file myenvfile -exec sh -c 'echo $E'
> a
```

For more information about env-files features see [github.com/joho/godotenv](https://github.com/joho/godotenv) documentation, which goenvtemplator internally uses as a library.

## Using Templates
Templates use Golang [text/template](http://golang.org/pkg/text/template/)
and [Sprig](https://github.com/Masterminds/sprig) library.

### Built-in functions
There are a few built in functions as well:
  * `{{ required  "Error message" Value }}` - Raises an error if `Value` is nil or it is equal to empty string. For example `{{ required "TIMEOUT_MS must be set!" (env "TIMEOUT_MS) }}`
  * `{{ range $key, $value := envall }}{{ $key }}={{ $value }};{{ end }}` - Loop over every environment variable.

### Nested Go templates
If you have nested Go templates there is problem with escaping. To resolve this problem you can define different
delimiters of template tags using `-delim-left` and `-delim-right` command line arguments.

Example of templating with `[[ ]]` instead of `{{ }}` delimiters
```bash
goenvtemplator2 -template /foo/bar.tmpl:/bar/foo.conf -delim-left [[ -delim-right ]]
```

## Development
Make sure you pass the CI. If the output of the GitHub actions is not descriptive enough,
or you can check it before submitting the code, use the `make test` and `make lint`.
