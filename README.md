# UVBOX

## Description
`uvbox` is a tool heavily inspired on the wonderful [pyapp](https://github.com/ofek/pyapp) idea, with opinionated flows & features, and focus on [uv](https://github.com/astral-sh/uv) benefits.

- 🌍 **`Cross-platform`**: out-of-the-box, without any licensing issue (MacOS).
- 🔥 **`Fast`**: Generate binaries for Windows, Linux and MacOS all in only few seconds.
- 📦 **`Packaged`**: Just add the dependency in your `pyproject.toml` and generate your binaries!
- 🐧 **`SOON`**: Generate debian and rpm packages at the same time

## Context
Python applications distribution is a pain 😡, and most of the time, you will ask your user:

#### With a virtual-environment 💬 (portable)
1. Please ensure you have at least `python3.X` installed
3. Please clone the repository
4. Please create a virtual environment
4. Please install the application dependencies

#### Or with a tool manager 💬 (portable)
1. Please ensure you have at least `python3.X` installed
2. Please ensure you have `pipx` (or something similar) installed
3. Please installed the application

#### Or with a compiled binary 💬 (not portable)
You will use something like [pyinstaller](https://pyinstaller.org/en/stable/). Here your user will just have to get and run your binary. Awesome! But the binary generation is almost always hard.
And contrary to the two last ways, cross-compilation of the binaries is usually very slow and painful, with a lot of dependencies required.
You can add on top of that licences issues, for example in order to compile MacOS binaries from Linux or Windows.

### So...
`uvbox` is trying to merge the best of these three solutions, in a fast and portable way:
- Generate a single binary for every platforms
- The binary will, at the first run, create its own python virtual environment and install the application
- Does not depends on the system python

Additionaly, we also provides extra feature, like auto-updates or dependencies freezing!

> <u>Note:</u> While [pyinstaller](https://pyinstaller.org/en/stable/) is an extremely powerful solution, here we want to only handle simple use-cases, and the simpliest manner. This project has no ambition to replace every features of pyinstaller.

## How?
The solution is based on [go](https://go.dev/).
This bring two huge added value:
- 🔥 The compilation is extremely **`fast`** (~1second). Rust-based [pyapp](https://github.com/ofek/pyapp) solution, in comparison, takes around ~30seconds on a M3 Pro MacBook Pro.
- 🚀 **`Cross-compilation`** **out-of-the-box**, without any system dependencies. \
This is possible because golang allows to generate binaries for every platforms, just by changing <b>GOOS</b> and <b>GOARCH</b> environment variables. Other solutions based on other languages like Rust most of the time requires a lot of system dependencies, or heavy container based setup.
- 🌍 Again on the cross-compilation: you can **`generate binaries for MacOS targets`**, from any platforms, <b>without any licences issues</b>. Again, thanks to [go](https://go.dev/).
- 📦 **`Easy setup and quickstart experience`**. The tool is available as Pypi package that you can just add as a development dependency inside your Python project. Even go will not be a requirements because it is also fetched from your pypi dependencies.

The binary will embed in itself its own [uv](https://github.com/astral-sh/uv) installation. <br>
At the very first run, [uv](https://github.com/astral-sh/uv) will be extracted inside a dedicated folder, and then installs your application, in a fully isolated manner.<br>
If python is missing on your system, it will automatically be downloaded in few tens of megabytes. <br>
You will be able to easily choose the python version of your choice.

You will also be able to easily configure your application to only target your own pypi mirror.

## Auto-update

One of the feature we wanted to provide out-of-the-box was updates.
Two kind of update strategy are possible:
- Run the command `<binary> self update` to update the application to the latest available version
- Configure auto-updates: the binary will, before every command, check for an update (This may slow response times).

You can also configure the way to compute the latest available version by providing your own URL of a remote text file containing the latest version tag. Allowing you to easily provide a fallback mechanism.

## How to use ?
Please have a look at the [HOWTOUSE]() page.

## How to contribute ?
Please have a look at the [CONTRIBUTING]() page.

## Licence
Please have a look at the [LICENCE]() page.
