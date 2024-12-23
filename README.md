# gopher
[![Release](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml)

**Gopher** is a Go project management CLI tool.
<br>
[Usage](#Usage) <span>&nbsp;•&nbsp;</span> [Preferences](#Preferences) <span>&nbsp;•&nbsp;</span> [Install](#Install)

## Usage
<details>
    <summary><code>gopher add</code></summary>

    > gopher help add
    $ Add dependencies to current module and install them.
    $ 
    $ When a full package name isn't provided gopher add will do a search on pkg.go.dev for matching packages. The number of results returned on this search can be adjusted with gopher config.
    $ 
    $ example: gopher add rsc.io/quote
    $ example: gopher add gofiber
</details>
<details>
    <summary><code>gopher build</code></summary>

    > gopher help build
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

    > gopher help help
    $ A Go project manager
    $ 
    $ usage: gopher [COMMAND] [...ARGS]
    $ 
    $ commands:
    $     add		add dependencies to current module and install them
    $     build	compile packages and dependencies
    $     help	this
    $     new		create new go module
    $     run		compile and run Go program
    $     test	run Go test packages
    $     tidy	add missing and remove unused modules
    $     version	print Go version
    $ 
    $ see gopher help [COMMAND] for more information about a specific command
</details>
<details>
    <summary><code>gopher new</code></summary>

    > gopher help new
    $ create new go module
    $ 
    $ example: gopher new github.com/user/mymodule
</details>
<details>
    <summary><code>gopher run</code></summary>

    > gopher help run
    $ compile and run Go program
    $ 
    $ gopher run should be executed at the root of your module and will expect the entry point of your program to be main.go
    $ 
    $ usage: gopher run [...ARGS]
    $ 
    $ example: gopher run
</details>
<details>
    <summary><code>gopher test</code></summary>

    > gopher help test
    $ run _test.go files
    $ 
    $ example: gopher test
</details>
<details>
    <summary><code>gopher tidy</code></summary>

    > gopher help tidy
    $ add missing and remove unused modules
    $ 
    $ example: gopher tidy
</details>
<details>
    <summary><code>gopher version</code></summary>

    > gopher help version
    $ print Go version
    $ 
    $ example: gopher version
</details>

## Preferences
On your initial call a **Preferences** file will be created at <code>~/.config/gopher</code>. Here you can customize aspects of the CLI to your liking.

    # The maximum number of modules returned on a 'gopher add' call
    PkgQueryLimit=10
    # List of architectures to target when running 'gopher build -x'
    OpArchPairs=windows,amd64,windows,arm64,linux,amd64,linux,arm64,darwin,amd64,darwin,arm64
    # Enables stylistic terminal output when an error is printed via 'gopher run'
    PrettyPrint=true
    # Number of lines printed to the terminal before and after the error line
    # Only takes effect when 'PrettyPrint' is set to 'true'
    PrettyPrintPreviewLines=3

## Install
Download pre-built binary for your system here [Releases](https://github.com/CoreyRobinsonDev/gopher/releases).

### Compiling from Source
- Clone this repository
```bash
git clone https://github.com/CoreyRobinsonDev/gopher.git
```
- Create **gopher** binary
```bash
cd gopher
go build
```
- Move binary to <code>/usr/local/bin</code> to call it from anywhere in the terminal
```bash
sudo mv ./gopher /usr/local/bin
```
- Confirm that the program was built successfully
```bash
gopher
```
    $ A Go project manager
    $ 
    $ usage: gopher [COMMAND] [...ARGS]
    $ 
    $ commands:
    $     add		add dependencies to current module and install them
    $     build	compile packages and dependencies
    $     help	this
    $     new		create new go module
    $     run		compile and run Go program
    $     test	run Go test packages
    $     tidy	add missing and remove unused modules
    $     version	print Go version
    $ 
    $ see gopher help [COMMAND] for more information about a specific command

## License
[Apache 2.0 License](./LICENSE)
