# How to use
Here you will find how to generate your binaries.

## Add uvbox to your python project

```console
# With uv
uv add --dev uvbox

# With poetry
poetry add --dev uvbox

# With pip
pip install uvbox
```

## Just run it

Just run the tool by giving your configuration file:
```console
❯ uvbox pypi --config eve.toml

██    ██ ██    ██ ██████   ██████  ██   ██
██    ██ ██    ██ ██   ██ ██    ██  ██ ██
██    ██ ██    ██ ██████  ██    ██   ███
██    ██  ██  ██  ██   ██ ██    ██  ██ ██
 ██████    ████   ██████   ██████  ██   ██

 SUCCESS  DARWIN/AMD64 👉 eve-x86_64-apple-darwin.tar.gz
 SUCCESS  DARWIN/ARM64 👉 eve-aarch64-apple-darwin.tar.gz
 SUCCESS  LINUX/AMD64 👉 eve-x86_64-unknown-linux-gnu.tar.gz
 SUCCESS  LINUX/ARM64 👉 eve-aarch64-unknown-linux-gnu.tar.gz
 SUCCESS  WINDOWS/AMD64 👉 eve-x86_64-pc-windows-msvc.zip
 SUCCESS  Available at: ./boxes
```

# Configuration file
First you will need to create a configuration file.

Some examples are [available here](https://github.com/AmadeusITGroup/uvbox/tree/master/examples).

The file is a `*.toml` file with the name of your choice, containing the following elements:

## `package`
The `package` required key allows you to configure which package will be installed at runtime, with specific version/updates if needed.

```toml
[package]
# Package to install.
name = "my-wonderful-package"
# Script to run.
# See: https://packaging.python.org/en/latest/guides/writing-pyproject-toml/#creating-executable-scripts
script = "main" 
```

### `package.version`
The `package.version`, required key allows you to configure for which version of the package the binary will ensure the installation.

```toml
[package]
...

[package.version]
# If set to true, the before each command, the binary will check and perform updates if possible 
auto-update = true
# If 'static' is set, the package will be initially installed with this specific version.
# The user will still be able to update/auto-update if needed
static = "1.2.3"
# You can specify a dynamic remote version file.
# The binary will download the file and reach the version from it.
# This is very useful in order to set-up fallback mechanisms.
dynamic = "https://www.myserver.example/MyWonderfulApp/version.txt
```

## `package.constraints`
The `package.constraints` key allows you to optionally use lockfile when installing your application. This useful because installing package regularly from pypi repositories does not provide dependencies freezing. When installing your application, a `requirements-txt` file will be used to constraint the installation. 

```toml
[package]
...

[package.constraints]
dynamic = "https://www.myserver.example/MyWonderfulApp/<VERSION>/version.txt"
```

<u>Note:</u> You can template the URL with the placeholder `<VERSION>`. `<VERSION>` will be matched and automatically replaced by the package version going to be installed.

# `certificates`
You may want to bundle a certificates bundle inside the binary in order to properly download your dependencies at runtime.
This is useful in company environments where you're behind company firewall.

By specifying a bundle to embed, the binary will automatically forward both `REQUESTS_CA_BUNDLE` and `SSL_CERT_FILE` environment variables to it.
```toml
[package]
...

[certificates]
path = "my-ca-bundle.crt"
```

<u>Note:</u> the path to the ca-bundle is relative to the working directory of the compilation.

# `uv`
`uv` key allows you to configure extra `uv` behaviour, with environment variables.
This will allows you to configure for example your pypi mirror, your python installation mirror, your indexes urls and so on.
[More informations here](https://docs.astral.sh/uv/configuration/environment/).

Example of an application with `uv` targetting an Artifactory, and use Python 3.13 for the installation:
```toml
[package]
name = "my-package"
script = "my-entrypoint"

...

[uv]
environment = [
    "UV_EXTRA_INDEX_URL=https://my.artifactory.com/artifactory/api/pypi/my-pypi-repository/simple",
    "UV_INDEX_URL=https://my.artifactory.com/artifactory/api/pypi/pypi-mirror/simple",
    "UV_PYTHON_INSTALL_MIRROR=https://my.github.remote/astral-sh/python-build-standalone/releases/download",
    "UV_PYTHON=3.13"
]
```
