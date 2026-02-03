<p align="center">
  <a href="https://www.kowabunga.cloud/?utm_source=github&utm_medium=logo" target="_blank">
    <picture>
      <source srcset="https://raw.githubusercontent.com/kowabunga-cloud/infographics/master/art/kowabunga-title-white.png" media="(prefers-color-scheme: dark)" />
      <source srcset="https://raw.githubusercontent.com/kowabunga-cloud/infographics/master/art/kowabunga-title-black.png" media="(prefers-color-scheme: light), (prefers-color-scheme: no-preference)" />
      <img src="https://raw.githubusercontent.com/kowabunga-cloud/infographics/master/art/kowabunga-title-black.png" alt="Kowabunga" width="800">
    </picture>
  </a>
</p>

# About

This is **Kiwi** Kowabunga agent, for SD-WAN nodes. It provides various network services like routing, firewall, DHCP, DNS, VPN, IPSec peering (with active-passive failover).

[![License: Apache License, Version 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://spdx.org/licenses/Apache-2.0.html)
[![Build Status](https://github.com/kowabunga-cloud/kiwi/actions/workflows/ci.yml/badge.svg)](https://github.com/kowabunga-cloud/kiwi/actions/workflows/ci.yml)
[![GoSec Status](https://github.com/kowabunga-cloud/kiwi/actions/workflows/sec.yml/badge.svg)](https://github.com/kowabunga-cloud/kiwi/actions/workflows/sec.yml)
[![GovulnCheck Status](https://github.com/kowabunga-cloud/kiwi/actions/workflows/vuln.yml/badge.svg)](https://github.com/kowabunga-cloud/kiwi/actions/workflows/vuln.yml)
[![Coverage Status](https://codecov.io/gh/kowabunga-cloud/kiwi/branch/master/graph/badge.svg)](https://codecov.io/gh/kowabunga-cloud/kiwi)
[![GoReport](https://goreportcard.com/badge/github.com/kowabunga-cloud/kiwi)](https://goreportcard.com/report/github.com/kowabunga-cloud/kiwi)
[![GoCode](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/kowabunga-cloud/kiwi)
[![time tracker](https://wakatime.com/badge/github/kowabunga-cloud/kiwi.svg)](https://wakatime.com/badge/github/kowabunga-cloud/kiwi)
![Code lines](https://sloc.xyz/github/kowabunga-cloud/kiwi/?category=code)
![Comments](https://sloc.xyz/github/kowabunga-cloud/kiwi/?category=comments)
![COCOMO](https://sloc.xyz/github/kowabunga-cloud/kiwiw/?category=cocomo&avg-wage=100000)

## Current Releases

| Project            | Release Badge                                                                                       |
|--------------------|-----------------------------------------------------------------------------------------------------|
| **Kiwi**           | [![Kiwi Release](https://img.shields.io/github/v/release/kowabunga-cloud/kiwi)](https://github.com/kowabunga-cloud/kiwi/releases) |

## Development Guidelines

Kiwi development relies on [pre-commit hooks](http://www.pre-commit.com/) to ensure proper commits.

Follow installation instructions [here](https://pre-commit.com/#install).

Local per-repository installation can be done through:

```sh
$ pre-commit install --install-hooks
```

And system-wide global installation, through:

```sh
$ git config --global init.templateDir ~/.git-template
$ pre-commit init-templatedir ~/.git-template
```

## Versioning

Versioning generally follows [Semantic Versioning](https://semver.org/).

## Authors

Kiwi is maintained by [Kowabunga maintainers](https://github.com/orgs/kowabunga-cloud/teams/maintainers).

## License

Licensed under [Apache License, Version 2.0](https://opensource.org/license/apache-2-0), see [`LICENSE`](LICENSE).
