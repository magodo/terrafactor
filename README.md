<p align="center">
  <a href="https://opensource.org/licenses/MIT"
    ><img
      src="https://img.shields.io/badge/License-MIT-yellow.svg"
      alt="License: MIT"
  /></a>
  <a href="https://github.com/magodo/terrafactor/actions"
    ><img
      src="https://img.shields.io/github/workflow/status/magodo/terrafactor/Go?label=workflow&style=flat-square"
      alt="GitHub Actions workflow status"
  /></a>
  <a href="https://twitter.com/magodo_"
    ><img
      src="https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Ftwitter.com%2Fmagodo_"
      alt="Follow me on Twitter"
  /></a>
</p>

🚧 This is very much WIP, do not use in production. 🚧

# Terrafactor

Terrafactor is a refactor tool for Terraform configurations.

## Install

Before Go v1.16:

```bash
$ go get github.com/magodo/terrafactor
```

Since Go v1.16:

```bash
$ go install github.com/magodo/terrafactor
```

## TODO

- [ ] Input Variable: Rename name (**cross modules**)
- [ ] Input Variable: Rename attribute for complex type (**cross modules**)
- [ ] Output Variable: Rename name (**cross modules**)
- [ ] Output Variable: Rename attribute for complex type (**cross modules**)
- [ ] Local Variable: Rename name
- [ ] Local Variable: Rename attribute for complex type
- [ ] Data/Managed Resource: Rename attribute/block name in complex expression (e.g. covering `Traversal` component of a `RelativeTraversalExpr`, `FunctionCallExpr`, etc)
