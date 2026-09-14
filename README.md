# yaml2yaml

`yaml2yaml` is a YAML preprocessor that transforms YAML documents containing
templating directives into plain YAML. It reads a YAML file (or standard input),
expands definitions, applies templates, evaluates conditional and loop
directives, resolves inline expressions, and writes the resulting plain YAML to
a file (or standard output).

## Installation

```
go install github.com/willoma/yaml2yaml/cmd/yaml2yaml@latest
```

## Usage

```
yaml2yaml [input] [output]
```

- With no arguments, `yaml2yaml` reads from standard input and writes to
  standard output.
- With one argument, it reads from the given file and writes to standard output.
- With two arguments, it reads from the first file and writes to the second.

```
cat example.yaml | yaml2yaml > output.yaml
yaml2yaml example.yaml output.yaml
```

The output is always re-emitted with 2-space indentation.

## Directives

Directives are written as YAML tags on mapping keys. yaml2yaml recognizes four
special tags: `!define`, `!apply`, `!if`, and `!for`.

### `!define` — define a reusable value

Binds a name to a value (scalar, sequence, or mapping). Definitions are scoped
to the mapping that contains them and are inherited by child mappings. A name
can only be defined once within a scope.

```yaml
!define majority: 18
!define year: 2026
```

A definition may itself contain further `!define` directives, building nested
templates:

```yaml
!define man:
  !define is: ugly
  !define has: smell
```

### `!apply` — apply a definition

Expands a previously defined value inline. When the applied value is a mapping,
the body of the `!apply` entry provides arguments that become definitions
visible inside the template:

```yaml
alice:
  !apply looks: f        # "f" is bound to the name "_"
  !apply adult:
    age: 34              # "age" is bound inside the "adult" template
```

When the argument is a single scalar or sequence, it is bound to the special
name `_`. When the argument is a mapping, each key becomes a definition.

### `!if` — conditional inclusion

Includes its body only when the expression in the key evaluates to a boolean
`true`. The expression uses the same syntax as inline expressions (see below).

```yaml
!define looks:
  !if _ == "f":
    !define is: strong
    !define has: boobies
  !if _ == "m": !apply man
```

### `!for` — loop over a sequence

Iterates over a previously defined sequence. For each item, the special names
`_` (the current item) and `i` (the zero-based index) are defined inside the
loop body. The body may be a mapping or a sequence.

```yaml
!define colors:
  - base: red
    variant: light
  - base: blue
    variant: dark

!for colors:
  ?_.variant + _.base:
    rank: ?i
    base_color: ?_.base
```

## Expressions

Any scalar value prefixed with `?` is evaluated as an expression. Expressions
are parsed using Go expression syntax (via `go/parser`) and resolved against the
current definitions.

```yaml
born in: ?year - age
adult for: ?age - majority + " years"
beautiful: ?_ == "f"
appearance: ?is + " with " + has
```

### Identifiers

A bare identifier resolves to a definition by name, walking up the scope chain
to parent mappings:

```yaml
- ?year        # the "year" definition
- ?majority    # the "majority" definition
```

The special name `_` refers to the current argument or loop item, and `i`
refers to the current loop index.

### Field access

Selector expressions (`.`) access fields of a mapping value:

```yaml
?_.base        # the "base" field of the current item
?_.variant     # the "variant" field of the current item
```

### Literals

Expressions support integer, float, and string literals using Go syntax:

```yaml
?age + 1
?1.5 * weight
?"hello " + name
```

### Operators

The following operators are supported:

| Category    | Operators                                  |
|-------------|--------------------------------------------|
| Arithmetic  | `+` `-` `*` `/` `%`                        |
| Comparison  | `==` `!=` `<` `>` `<=` `>=`                |
| Logical     | `&&` `\|\|`                                |

Arithmetic and comparison operators perform automatic type coercion between
integers, floats, strings, booleans, and byte slices where it makes sense. For
example, `+` concatenates strings and sequences, `*` repeats strings and
sequences, and numeric strings are parsed when combined with numbers.

## Example

Given the input `example.yaml`:

```yaml
!define majority: 18
!define year: 2026
!define man:
  !define is: ugly
  !define has: smell

!define looks:
  !if _ == "f":
    !define is: strong
    !define has: boobies
  !if _ == "m": !apply man
  beautiful: ?_ == "f"
  appearance: ?is + " with " + has

!define adult:
  born in: ?year - age
  adult for: ?age - majority + " years"

alice:
  !apply looks: f
  job: CEO
  !apply adult:
    age: 34
bob:
  !apply looks: m
  !apply adult:
    age: 28
  job: assistant

!define colors:
  - base: red
    variant: light
  - base: blue
    variant: dark
  - base: green
    variant: medium

!for colors:
  ?_.variant + _.base:
    rank: ?i
    base_color: ?_.base

colorslist:
  !for colors:
    - name: ?_.variant + _.base
      base_color: ?_.base
```

Running `yaml2yaml example.yaml` produces:

```yaml
alice:
  beautiful: true
  appearance: strong with boobies
  job: CEO
  born in: 1992
  adult for: 16 years
bob:
  beautiful: false
  appearance: ugly with smell
  born in: 1998
  adult for: 10 years
  job: assistant
lightred:
  rank: 0
  base_color: red
darkblue:
  rank: 1
  base_color: blue
mediumgreen:
  rank: 2
  base_color: green
colorslist:
- name: lightred
  base_color: red
- name: darkblue
  base_color: blue
- name: mediumgreen
  base_color: green
```

## Library usage

yaml2yaml can also be used as a Go library:

```go
import "github.com/willoma/yaml2yaml"

result, err := yaml2yaml.Bytes(input)
```

`yaml2yaml.Bytes` parses the input YAML, processes all directives and
expressions, and returns the resulting YAML bytes. A lower-level
`yaml2yaml.Process` function is also available for operating directly on
`yaml.Node` values.

## License

This project is licensed under the [MIT License](LICENSE).
