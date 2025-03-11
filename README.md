# wbox

CLI tool to manage weedbox projects.

## Installation

```shell
go install github.com/weedbox/wbox@latest
```

## Usage

### Initialize Weedbox project

```shell
mkdir myproject
cd myproject
wbox init myproject github.com/myuser/myproject
```

### Createa a new module in the Weedbox project

```shell
cd pkg
mkdir mymodule
wbox init-module MyModule
```
