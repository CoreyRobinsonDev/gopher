# gopher
[![Release](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml/badge.svg)](https://github.com/CoreyRobinsonDev/gopher/actions/workflows/release.yml)
[![Report](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/CoreyRobinsonDev/gopher)

**Gopher** is a Go project management CLI tool.
<br>
[Usage](#Usage) <span>&nbsp;•&nbsp;</span> [Config](#Config) <span>&nbsp;•&nbsp;</span> [Install](#Install)

## Usage
<details>
    <summary><code>gopher add</code></summary>
    <img alt="add command gif" src="https://vhs.charm.sh/vhs-VxKxN5my8JTDsSu5HPijo.gif"/>
    <img alt="add command example gif" src="https://vhs.charm.sh/vhs-6SPu40VY4S6egGciFxYl4U.gif"/>
</details>
<details>
    <summary><code>gopher build</code></summary>
    - use <b>--web</b> to compile the program for the browser
    <img alt="build command gif" src="https://vhs.charm.sh/vhs-6HaDjeb6NN2SbkCJOGR3ox.gif"/>
</details>
<details>
    <summary><code>gopher config</code></summary>
    <img alt="config command gif" src="https://vhs.charm.sh/vhs-6hQdFU1ZSXCZEd4yNb4mlB.gif"/>
</details>
<details>
    <summary><code>gopher new</code></summary>
    <img alt="new command gif" src="https://vhs.charm.sh/vhs-69YDFALfOTecVT1HmJjFHP.gif"/>
</details>
<details>
    <summary><code>gopher release</code></summary>
    - use <b>-c, --clean</b> to remove the 'dist' directory before build <br>
    <img alt="release command gif" src="https://vhs.charm.sh/vhs-15PT4fvtcigaoZzq8YNFM0.gif"/>
</details>
<details>
    <summary><code>gopher run</code></summary>
    - use <b>--web</b> to run the program in the browser <br>
    - use <b>-w, --watch</b> to live-reload your code on change <br>
    <img alt="run command gif" src="https://vhs.charm.sh/vhs-5aIRsVxYDlZpMRMSDGlWJP.gif"/>

</details>
<details>
    <summary><code>gopher test</code></summary>
    <img alt="test command gif" src="https://vhs.charm.sh/vhs-2tCXkm2NSVWSj6sNZ4JEC7.gif"/>
</details>
<details>
    <summary><code>gopher tidy</code></summary>
    <img alt="tidy command gif" src="https://vhs.charm.sh/vhs-2NJcaxNnzj9jf9g0nZseAU.gif"/>
</details>
<details>
    <summary><code>gopher version</code></summary>
    <img alt="version command gif" src="https://vhs.charm.sh/vhs-2mhDWhXegEYaUO6LWSuh2u.gif"/>
</details>

## Config
On your initial call a **settings.json** file will be created at <code>~/.config/gopher</code>. Here you can customize aspects of the CLI to your liking.
Default values are as shown:

```jsonc
{
    // enables stylized output for errors when running 'gopher run' and 'gopher build'
	"prettyPrint": true,
    // the number of lines around the the errored line to output when 'prettyPrint' is enabled
	"prettyPrintPreviewLines": 3,
    // the number of packages to query for when 'gopher add {package}' is ran
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

## License
[Apache 2.0 License](./LICENSE)
