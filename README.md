[![Build Status](https://travis-ci.org/unjello/belit.svg?branch=master)](https://travis-ci.org/unjello/belit)
[![Coverage Status](https://codecov.io/gh/unjello/belit/branch/master/graph/badge.svg)](https://codecov.io/gh/unjello/belit)
[![Go Report Card](https://goreportcard.com/badge/github.com/unjello/belit)](https://goreportcard.com/report/github.com/unjello/belit)
[![GoDoc](https://godoc.org/github.com/unjello/belit?status.svg)](https://godoc.org/github.com/unjello/belit)

# Bêlit?

Bêlit is a pirate queen in romantic relationship with Conan. It is also an extremely simple package manager for C/C++ that works with header-only libraries.

## Why Bêlit?

I found myself writing a lot of single-file C++ tidbits, that could use unit testing. Setting up build system in C++ is tedious, especially
if you do compile-run loop on a number of small, different, programs.

Enter Bêlit. Heavily inspired by `Go`. It uses non-intrusive tagging in C++ file:
```c++
#include /* github.com/catchorg/Catch2/single_include/ */ "catch.hpp"
```
where a comment between `#include` keyword and a header name, specifies where in remote repository (currently only GitHub is supported) a header can be found.
This comment is totally innocent, and invisible to the rest of your regular C++ ecosystem, so same source can easily work side-by-side with other generators like [Cmake](http://cmake.org) or package managers like [conan.io](http://conan.io) if you choose to use them.

- Bêlit `get` command will fetch headers from remote repository and cache them.
- Bêlit `run` command will compile the program, adjusting header search paths accordingly, and run it.

## Configuration

Bêlit honours `CXX`, `CC`, `CXXOPTS` and `CCOPTS` environment variables. You can also use configuration file in TOML, JSON or YAML formats. Bêlit searches for file named `config.toml`, `config.json` or `config.yml` in current `.` folder, `$HOME/.belit/` or `/etc/belit/`.

```yaml
cxx: /usr/local/bin/g++-7
cc: /usr/local/bin/gcc-7
cxxopts: -std=c++17 -O2
```

it also works with modeline embedded options, vim-style:

```c
int main() {
  return 0;
}
/* belit: cxx=g++-7 cxxopts="-O2 -std=c++17" */
```

## Example

### Simple example

Most basic example is compiling a program with no dependencies. In this case, `belit` is just a convenient shorthand for compilation without setting up build systems.
```bash
$ belit run example/simple.cpp
Hello C++ world
```

### Headers example

In a little more complicated example, we have a dependency on excellent unit test library: [Catch](https://github.com/catchorg/Catch2).

```batch
$ belit get example/catch.cpp

$ belit run example/catch.cpp
===============================================================================
All tests passed (5 assertions in 2 test cases)
```

## License

- Unlicensed (~Public Domain)

## Related Work

- [Go](https://golang.org) - heavily inspired by how packages work in Go
- [conan.io](https://conan.io) - proper package manager for C/C++