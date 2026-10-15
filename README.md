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
yaml2yaml [-i input] [-o output]
```

- With no flags, `yaml2yaml` reads from standard input and writes to
  standard output.
- `-input` (or `-i`) sets the file to read from.
- `-output` (or `-o`) sets the file to write to.

```
cat example/full.yaml | yaml2yaml > output.yaml
yaml2yaml -i example/full.yaml -o output.yaml
```

YAML anchors and aliases are resolved transparently. The output is always
re-emitted with 2-space indentation.

## Directives

Directives are written as YAML tags on mapping keys.
yaml2yaml recognizes the special tags documented hereafter.

When several entries of a mapping produce the same key, the new value is
merged into the existing one if it is a mapping (failing if the new value is
not itself a mapping) or a sequence, and fails otherwise.

### `!define` — define a reusable value

Binds a name to a value (scalar, sequence, or mapping). Definitions are scoped
to the mapping that contains them and are inherited by child mappings.
Defining an already-defined name in the same scope replaces the previous
definition.

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

When a construct produces no output of its own — typically because it only
contains `!define` directives — the definitions it created are propagated to
the enclosing scope. This applies to `!if` bodies, applied templates, and
included files as well as to regular mappings.

### `!apply` — apply a definition

Expands a previously defined value inline. When the applied value is a mapping,
the body of the `!apply` entry provides arguments that become definitions
visible inside the template:

```yaml
alice:
  !apply looks: f   # "f" is bound to the name "_"
  !apply adult:
    age: 34         # "age" is bound inside the "adult" template
```

When the argument is a single scalar or sequence, it is bound to the special
name `_`. When the argument is a mapping, each plain `key: value` entry
becomes a definition; directive entries are ignored.

`!apply` can also be used as a scalar value (`key: !apply name`), in which case
the value is replaced by the named definition.

Applying a name that is not defined is an error — unlike expressions, `!apply`
requires the definition to exist.

If the applied value produces no output — typically a template containing only
`!define` directives — its definitions are propagated to the scope containing
the `!apply` entry:

```yaml
!if _ == "m": !apply man   # "is" and "has" become defined in this scope
```

### `!expr` - expressions

Evaluates the value as an expression,
parsed using Go expression syntax (via `go/parser`)
and resolved against the current definitions.

```yaml
born in: !expr year - age
adult for: !expr age - majority + " years"
beautiful: !expr _ == "f"
appearance: !expr is + " with " + has
```

### `!fmt` — formatted string

Formats a string, replacing "`$(xxx)`" sections with the result
of expression `xxx`. The `$` character can be kept and not interpreted,
by escaping it with `\`.

```yaml
!fmt In a formatted string, \$(year) would be processed as "$(year)"
```

### `!for` — loop over a sequence

Iterates over a sequence. The key may be the name of a previously defined
sequence, an expression resolving to such a name, or an inline sequence. If a
name is not defined, the loop produces nothing and does not fail; if it is
defined but is not a sequence, an error is reported. If the body is a mapping,
the result will be a mapping, merging the result of the loop calculation.
If the body is a sequence, the result will be a sequence, joining all
calculated sequences. If the body is a scalar, it returns an error.

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
  !expr variant + base:
    rank: !expr i
    base_color: !expr base
```

### `!if` — conditional inclusion

Includes its body only when the expression in the key is truthy.
The expression uses the same syntax as with the `!expr` directive.
Truthy values are `true`, positive numbers,
and non-empty strings, sequences, mappings and binary data;
everything else — `false`, zero and negative numbers,
empty strings, sequences, mappings and binary data — is falsy.
An identifier that is not defined produces no value,
so the `!if` entry is simply skipped instead of failing.

```yaml
!define looks:
  !if _ == "f":
    !define is: strong
    !define has: boobies
  !if _ == "m": !apply man
```

A body that produces no output — typically one containing only `!define`
directives — passes its definitions to the enclosing scope, as in the example
above where `is` and `has` become defined in the `looks` template.

### `!include` — include another file

Loads another YAML file, processes it, and inlines the result in place of the
directive. A relative path is resolved against the directory of the file
containing the directive — or against the current working directory when
reading from standard input.
An `!include` entry does not use a value; any value provided is ignored.

```yaml
!include common/definitions.yaml:
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

Values evaluated with the `!expr`, `!fmt` and `!if` directives
follow the logic described hereafter.

### Identifiers

A bare identifier resolves to a definition by name, walking up the scope chain
to parent mappings:

```yaml
- !expr year       # the "year" definition
- !expr majority   # the "majority" definition
```

`true` and `false` are boolean literals. An identifier that is not defined
produces no value: the expression is considered empty and the entry containing
it is omitted from the output instead of being rendered (in an `!if`
condition, this simply skips the entry).

### Field access

Selector expressions (`.`) access fields of a mapping value:

```yaml
- !expr color.base      # the "base" field of the current item
- !expr color.variant   # the "variant" field of the current item
```

Accessing a field that does not exist produces no value. Accessing a field of
a value that is not a mapping is an error.

### Element access

Index expressions (`[]`) access elements of a sequence:

```yaml
- !expr colors[0]        # the first element
- !expr colors[0, 2]     # a sequence of the first and third elements
- !expr colors[1:3]      # a slice containing elements 1 and 2
- !expr colors[3].base   # indexing can be chained with field access
```

Indexes and slice bounds must be integers and are bounds-checked; out-of-range
access is an error. In a slice, both bounds are required and the start must be
strictly before the end.

### Literals

Expressions support integer, float, character, and string literals using Go
syntax, plus the booleans `true` and `false`:

```yaml
- !expr age + 1
- !expr 1.5 * weight
- !expr "hello " + name
```

### Operators

The following operators are supported:

| Category   | Operators                        |
| ---------- | -------------------------------- |
| Arithmetic | `+`, `-`, `*`, `/`, `%`          |
| Comparison | `==`, `!=`, `<`, `>`, `<=`, `>=` |
| Logical    | `&&`, `\|\|`, `!`                |

Arithmetic:

- `a + b` adds numbers, concatenates strings, sequences, and binary data, and
  merges mappings. Adding a non-sequence value to a sequence appends it.
- `a - b` and `a / b` subtract and divide integers and floats.
- `a * b` multiplies integers and floats. When `b` is an integer, it instead
  repeats a string, sequence, or binary `a` (`!expr "ab" * 3` produces `"ababab"`).
- `a % b` is the remainder of an integer division.

Comparison:

- `a == b`, `a != b`, `a < b`, `a > b`, `a <= b`, `a >= b` compare integers,
  floats, and strings.
- Sequences and mappings are compared by length first, then element by
  element; binary data is compared byte by byte.
- Booleans only support the equality-style comparisons (`==`, `!=`, `<=`,
  `>=`).

Logical:

- `a && b`, `a || b`, and `!a` operate on truthiness (see `!if` above).

Arithmetic and comparison operators perform automatic type coercion: operands
are tried in order as integers, floats, booleans, then strings, so numeric
strings are parsed when combined with numbers (`!expr "5" * 3` produces `15`, not
`"555"`).

Parentheses can be used to group expressions.

## Example

Given the input `example/full.yaml`, which includes
`example/full.include.yaml`:

```yaml
# example/full.include.yaml
!define majority: 18
!define year: 2026
!define man:
  !define is: ugly
  !define has: smell
```

```yaml
# example/full.yaml
!include full.include.yaml:

!define looks:
  !if _ == "f":
    !define is: strong
    !define has: boobies
  !if _ == "m": !apply man
  beautiful: !expr _ == "f"
  appearance: !expr is + " with " + has

!define adult:
  born in: !expr year - age
  adult for: !expr age - majority + " years"
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
  - base: yellow
    variant: pastel
  - base: grey
    variant: light

!for colors:
  !expr variant + base:
    rank: !expr i
    base_color: !expr base

colorslist:
  !for colors:
    - name: !expr variant + base
      base_color: !expr base

somecolors: !expr colors[1:3]

fourthcolor: !expr colors[3].base

!if colors:
  has_colors: true
```

Running `yaml2yaml -i example/full.yaml` produces:

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
pastelyellow:
  rank: 3
  base_color: yellow
lightgrey:
  rank: 4
  base_color: grey
colorslist:
  - name: lightred
    base_color: red
  - name: darkblue
    base_color: blue
  - name: mediumgreen
    base_color: green
  - name: pastelyellow
    base_color: yellow
  - name: lightgrey
    base_color: grey
somecolors:
  - base: blue
    variant: dark
  - base: green
    variant: medium
fourthcolor: yellow
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
  bytes. `!include` paths are resolved against the current working directory.
- `yaml2yaml.ProcessFile(path string) ([]byte, error)` does the same starting
  from a file path; `!include` paths are resolved relative to that file's
  directory.
- `yaml2yaml.Process(basePath string, src *yaml.Node) (*yaml.Node, error)` is
  the lower-level entry point, operating directly on `yaml.Node` values; the
  `basePath` argument is the directory against which relative `!include`
  paths are resolved.

## License

This project is licensed under the [MIT License](LICENSE).
