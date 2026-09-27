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

Directives are written as YAML tags on mapping keys. yaml2yaml recognizes five
special tags: `!define`, `!apply`, `!if`, `!for`, and `!include`.

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

`!apply` can also be used as a scalar value (`key: !apply name`), in which case
the value is replaced by the named definition.

Applying a name that is not defined is an error — unlike expressions, `!apply`
requires the definition to exist.

### `!if` — conditional inclusion

Includes its body only when the expression in the key is truthy. The
expression uses the same syntax as inline expressions (see below), but is not
restricted to booleans: `false`, `0`, `0.0`, `""`, empty sequences, empty
mappings, and `nil` are falsy; any other value is truthy. An identifier that
is not defined resolves to `nil`, so the condition simply evaluates to false
instead of failing.

```yaml
!define looks:
  !if _ == "f":
    !define is: strong
    !define has: boobies
  !if _ == "m": !apply man
```

### `!for` — loop over a sequence

Iterates over a previously defined sequence. If the name is not defined, the
loop produces nothing and does not fail; if it is defined but is not a
sequence, an error is reported. The body must be a mapping or a sequence.

For each item, the zero-based index is bound to the name `i`. When the item
is a mapping, each of its keys becomes a definition inside the loop body;
when the item is a scalar or a sequence, the item itself is bound to the
special name `_`.

```yaml
!define colors:
  - base: red
    variant: light
  - base: blue
    variant: dark

!for colors:
  ?variant + base:
    rank: ?i
    base_color: ?base
```

### `!include` — include another file

Loads another YAML file, processes it, and inlines the result in place of the
directive. The path is relative to the current working directory, and the key
may be an expression. An `!include` entry must not have a value.

```yaml
!include common/definitions.yaml:
!include ?dir + "/" + name + ".yaml":
```

The included file is processed with access to the definitions in scope at the
include point. An included file containing only `!define` directives passes
its definitions to the including scope, making it a convenient way to share
common definitions.

`!include` can also be used as a scalar value, in which case it is replaced by
the processed content of the file:

```yaml
fragment: !include fragment.yaml
```

If the path is a directory, all the files it contains are included
recursively, in filename order, and merged together: mappings are merged key
by key, sequences are concatenated, and scalars from later files replace
earlier ones. Merging incompatible kinds is an error.

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

An identifier that is not defined resolves to `nil` without failing. `nil` is
falsy, so `!if` conditions and the `&&`/`||` operators treat undefined values
as `false`. In comparisons, `nil` never equals anything: `?undefined == x` is
always `false` — even when `x` is itself undefined — while `?x != undefined`
is `true` whenever `x` is defined. Used alone, a `nil` value is rendered as
`<nil>`.

The special name `_` refers to the current argument or loop item, and `i`
refers to the current loop index.

### Field access

Selector expressions (`.`) access fields of a mapping value:

```yaml
- ?color.base        # the "base" field of the current item
- ?color.variant     # the "variant" field of the current item
```

Accessing a field of a value that is not a mapping (including an undefined
value) is an error.

### Literals

Expressions support integer, float, and string literals using Go syntax:

```yaml
?age + 1
?1.5 * weight
?"hello " + name
```

### Operators

The following operators are supported:

| Category   | Operators                   |
| ---------- | --------------------------- |
| Arithmetic | `+` `-` `*` `/` `%`         |
| Comparison | `==` `!=` `<` `>` `<=` `>=` |
| Logical    | `&&` `\|\|`                 |

Arithmetic and comparison operators perform automatic type coercion between
integers, floats, strings, booleans, and byte slices where it makes sense. For
example, `+` concatenates strings and sequences, `*` repeats strings and
sequences, and numeric strings are parsed when combined with numbers.

## Example

Given the input `example.yaml`, which includes `example.include.yaml`:

```yaml
# example.include.yaml
!define majority: 18
!define year: 2026
!define man:
  !define is: ugly
  !define has: smell
```

```yaml
# example.yaml
!include example.include.yaml:

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
  !if wonderful:
    person: wonderful
  !if wonderful == false:
    person: mediocre

alice:
  !apply looks: f
  job: CEO
  !apply adult:
    age: 34
    wonderful: true
bob:
  !apply looks: m
  !apply adult:
    age: 28
    wonderful: false
  job: assistant

!define colors:
  - base: red
    variant: light
  - base: blue
    variant: dark
  - base: green
    variant: medium

!for colors:
  ?variant + base:
    rank: ?i
    base_color: ?base

colorslist:
  !for colors:
    - name: ?variant + base
      base_color: ?base

!if colors:
  has_colors: true
```

Running `yaml2yaml example.yaml` produces:

```yaml
alice:
  beautiful: true
  appearance: strong with boobies
  job: CEO
  born in: 1992
  adult for: 16 years
  person: wonderful
bob:
  beautiful: false
  appearance: ugly with smell
  born in: 1998
  adult for: 10 years
  person: mediocre
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
has_colors: true
```

## Library usage

yaml2yaml can also be used as a Go library:

```go
import "github.com/willoma/yaml2yaml"

result, err := yaml2yaml.ProcessFile("input.yaml")
```

Three entry points are available:

- `yaml2yaml.ProcessBytes(src []byte) ([]byte, error)` parses the input YAML,
  processes all directives and expressions, and returns the resulting YAML
  bytes.
- `yaml2yaml.ProcessFile(path string) ([]byte, error)` does the same starting
  from a file path.
- `yaml2yaml.Process(src *yaml.Node) (*yaml.Node, error)` is the lower-level
  entry point, operating directly on `yaml.Node` values.

## License

This project is licensed under the [MIT License](LICENSE).
