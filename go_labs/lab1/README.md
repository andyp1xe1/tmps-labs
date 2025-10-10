# Lab 1 - Text Search Tool

A command-line text search tool demonstrating SOLID principles (S, D, L) through clean interface design and dependency injection.

## SOLID Principles Overview

The **SOLID** principles are five design principles that make software designs more understandable, flexible, and maintainable:

- **S** - Single Responsibility Principle: Each class should have one reason to change
- **O** - Open/Closed Principle: Open for extension, closed for modification  
- **L** - Liskov Substitution Principle: Objects should be replaceable with instances of their subtypes
- **I** - Interface Segregation Principle: Many specific interfaces are better than one general-purpose interface
- **D** - Dependency Inversion Principle: Depend on abstractions, not concretions

This lab focuses on **S**, **D**, and **L** principles.

## Usage

```bash
go run . -e <engine> -q <query> -f <format> -p <path>
```

### Flags

- `-e` - Search engine: `literal`, `regex`, `fuzzy`
- `-q` - Search query
- `-f` - Output format: `plain`, `json`
- `-p` - File path to search

### Examples

```bash
# Literal search with plain output
go run . -e literal -q "world" -f plain -p test.txt

# Regex search with JSON output
go run . -e regex -q "^[A-Z].*ing" -f json -p test.txt

# Fuzzy search
go run . -e fuzzy -q "hlwrd" -f plain -p test.txt
```

## Testing

```bash
go test -v
```

## Search Types

- **Literal**: Substring matching
- **Regex**: Regular expression patterns
- **Fuzzy**: Character sequence matching (order matters, gaps allowed)

## Architecture & SOLID Principles

### Interfaces Defined

The application defines minimal interfaces where they're needed:

```go
// In runner.go - where SearchEngine is used
type SearchEngine interface {
    Search(text, query string) bool
}

// In runner.go - where ResultWriter is used  
type ResultWriter interface {
    Write(results []SearchResult) error
}
```

### Single Responsibility Principle (S)

Each component has one clear responsibility:

```go
// LiteralSearch - only handles substring matching
type LiteralSearch struct{}
func (l *LiteralSearch) Search(text, query string) bool {
    return strings.Contains(text, query)
}

// PlainWriter - only handles plain text output
type PlainWriter struct{ output io.Writer }
func (p *PlainWriter) Write(results []SearchResult) error {
    // Format as plain text only
}
```

### Dependency Inversion Principle (D)

High-level modules depend on abstractions, not concrete implementations:

```go
// Runner depends on interfaces, not concrete types
type Runner struct {
    engine SearchEngine  // Interface, not *LiteralSearch
    reader io.Reader     // Standard interface
    writer ResultWriter // Interface, not *PlainWriter
}

func NewRunner(engine SearchEngine, reader io.Reader, writer ResultWriter) *Runner {
    return &Runner{engine: engine, reader: reader, writer: writer}
}
```

### Liskov Substitution Principle (L)

All implementations can be used interchangeably without breaking functionality:

```go
// Any SearchEngine implementation works the same way
engines := []SearchEngine{
    &LiteralSearch{},
    &RegexSearch{},
    &FuzzySearch{},
}

// Any ResultWriter implementation works the same way
writers := []ResultWriter{
    &PlainWriter{output: os.Stdout},
    &JSONWriter{output: os.Stdout},
}
```

### Search Engine Implementations

- **LiteralSearch**: Simple `strings.Contains()` for exact substring matching
- **RegexSearch**: Uses `regexp.MatchString()` for pattern-based searching with full regex support
- **FuzzySearch**: Subsequence matching algorithm allowing gaps between characters while preserving order

### Writer Implementations

- **PlainWriter**: Human-readable format `"lineNumber: content"` for terminal output
- **JSONWriter**: Structured format for programmatic consumption with `line_number` and `line` fields

## Conclusion

This implementation demonstrates how SOLID principles create modular, testable code. The `Runner` orchestrates components without knowing their concrete types, making the system easily extensible. New search algorithms or output formats can be added without modifying existing code, showcasing the power of interface-based design and dependency injection.
