# gopher
[![Release](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml)
[![Report](https://goreportcard.com/badge/github.com/CoreyRobinsonDev/gopher)](https://goreportcard.com/report/github.com/CoreyRobinsonDev/gopher)

Gopher is a Go project management CLI tool.

# Usage
<details>
    <summary><code>gopher add</code></summary>

    $ gopher help add
    $ Add dependencies to current module and install them.
    $ 
    $ When a full package name isn't provided gopher add will do a search on pkg.go.dev for matching packages. The number of results returned on this search can be adjusted with gopher config.
    $ 
    $ example: gopher add rsc.io/quote
    $ example: gopher add gofiber
</details>
<details>
    <summary><code>gopher build</code></summary>
$ gopher help build
$ compile packages and dependencies
$ 
$ gopher build should be executed at the root of your module and will expect the entry point of your program to be main.go
$ 
$ usage: gopher build [...ARGS]
$ 
$ arguments:
$     -x,--cross-platform		build binaries for seperate operating systems and cpu architectures speficied by your gopher configuration
$ 
$ example: gopher build
</details>
<details>
    <summary><code>gopher help</code></summary>
</details>
<details>
    <summary><code>gopher new</code></summary>
</details>
<details>
    <summary><code>gopher run</code></summary>
</details>
<details>
    <summary><code>gopher test</code></summary>
</details>
<details>
    <summary><code>gopher tidy</code></summary>
</details>
<details>
    <summary><code>gopher version</code></summary>
</details>
