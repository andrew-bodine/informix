# Informix

## Development
It is assumed that all `Informix` development is done on some machine that is
not the deployment target. For example, even though we want `Informix` running
on some edge device, we want to develop on our development machines.

#### Build & Test
In order to build a development version of `Informix`, run the following after
making whatever changes you are working on to the source:
```
$ make build
```

This will verify the source still builds correctly, and will proceed to run all
unit and integration tests with `Ginkgo`. A successful build here means you are
good to go for checking in your changes.
