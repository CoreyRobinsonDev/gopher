# gopher
[![Release](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml/badge.svg)](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml)
[![Report](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/CoreyRobinsonDev/gopher)

**Gopher** is a Go project management CLI tool.
<br>
[Usage](#Usage) <span>&nbsp;•&nbsp;</span> [Preferences](#Preferences) <span>&nbsp;•&nbsp;</span> [Install](#Install)

## Usage
<details>
    <summary><code>gopher add</code></summary>
    ![add](https://vhs.charm.sh/vhs-VxKxN5my8JTDsSu5HPijo.gif)
</details>
<details>
    <summary><code>gopher build</code></summary>
    ![build](https://vhs.charm.sh/vhs-4cqk1DmrECFWHnHE21kLyA.gif)
</details>
<details>
    <summary><code>gopher new</code></summary>
    ![new](https://vhs.charm.sh/vhs-69YDFALfOTecVT1HmJjFHP.gif)
</details>
<details>
    <summary><code>gopher run</code></summary>
    ![run](https://vhs.charm.sh/vhs-32VJdIqvYHoH8wMk6grmZW.gif)
</details>
<details>
    <summary><code>gopher test</code></summary>
    ![test](https://vhs.charm.sh/vhs-2tCXkm2NSVWSj6sNZ4JEC7.gif)
</details>
<details>
    <summary><code>gopher tidy</code></summary>
    ![tidy](https://vhs.charm.sh/vhs-2NJcaxNnzj9jf9g0nZseAU.gif)
</details>
<details>
    <summary><code>gopher version</code></summary>
    ![version](https://vhs.charm.sh/vhs-2mhDWhXegEYaUO6LWSuh2u.gif)
</details>

## Config
On your initial call a **settings.json** file will be created at <code>~/.config/gopher</code>. Here you can customize aspects of the CLI to your liking.
Default values are as shown:

```json
{
	"prettyPrint": true,
	"prettyPrintPreviewLines": 3,
	"pkgQueryLimit": 10
}
```

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
