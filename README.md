# 1lang Interpreter

[简体中文版本](README.ZH_CN.md)

This project is a custom interpreter written in Go, implementing the Lexer, Parser, AST (Abstract Syntax Tree), and Evaluator.

The project is inspired by Thorsten Ball's book "Writing An Interpreter In Go". Building on the original implementation, new keywords and feature modifications have been added to meet the characteristics of modern programming languages.

## Now Supported Features
- Variable assignment and reference
- Function definition and call
- Arithmetic operations
- Conditional statements
- Loop statements [While]
- Array operations
- Importing external modules
- Comments

## Current Issues
- Error reporting lacks line and column information

## Future Work
- Macro system

## Example Usage

You can learn this new language using a simple example:

```javascript
const filter = fn(arr, f) {
  let iter = fn(arr, acc) {
    if (len(arr) == 0) {
      acc;
    } else {
      let x = first(arr);
      let restArr = rest(arr);
      if (f(x)) {
        iter(restArr, push(acc, x));
      } else {
        iter(restArr, acc);
      }
    }
  };
  iter(arr, []);
};

const quick_sort = fn(arr) {
  if (len(arr) <= 1) {
    arr;
  } else {
    let pivot = first(arr);
    let restArr = rest(arr);
    let less = filter(restArr, fn(x) { x <= pivot });
    let greater = filter(restArr, fn(x) { x > pivot });
    concat(quick_sort(less), [pivot], quick_sort(greater));
  }
};

quick_sort([3, 6, 8, 10, 1, 2, 1]);
```

## REPL and File Execution Support

To run the interpreter in REPL mode or to execute a file, use the following command:

```bash
go run main.go -f <file_path> -t
```

`t` is an optional parameter that enables the time measurement of the execution process.

This project aims to provide a learning platform for building interpreters and understanding the intricacies of programming language design. Contributions and feedback are welcome!