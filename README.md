# Building a Treewalking Interpreter with Go

## Build Components:
- Lexer
- Parser
- Tree Representation
- Internal Object System
- Evaluator

## Features of the Programming Language:

- **C-like Syntax**
- **Variable Bindings with Type Inference**
- **Supported Data Types:**
  - Integers
  - Booleans
  - Strings
  - Pointers
- **Expressions:**
  - Arithmetic Expressions
  - `while` Loop
- **Built-in Functions**
- **First-Class Functions & Higher Order Functions**
- **Closures**
- **Data Structures:**
  - Arrays
  - Maps
- **Heap Memory**
- **Garbage Collection**

### Examples:

#### Variable Bindings with `let`
```go
let numbers = [1, 23, 42];
```

#### Map Key-Value Pairs
```go
let me = {"name": "Mahmut", "profession": "Student"};
```

#### Functions Bound to Variables
```go
let add = func(x, y) {
    x + y;
};
```

#### Higher-Order Functions
Functions can take other functions as arguments and return functions.
```go
let twice = func(f, x) {
    f(f(x));
};
```


