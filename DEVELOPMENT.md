# yaml2yaml development rules

This document outlines some development rules for the `yaml2yaml` project.

## Invariants

**No in-place mutation of nodes.** All processing operations must build new
`*yaml.Node` values. Original nodes (including definitions) must remain
unmodified — future features may need access to the originals.

## Example-based test

The base test for `yaml2yaml` is an example-based test that verifies
that the tool can convert a YAML file to a YAML file.

It reads the `example.yaml` file and converts it to a YAML file,
that must match the `example.output.yaml` file.
