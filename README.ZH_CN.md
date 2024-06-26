# 1lang 解释器

[English Version](README.md)

本项目是由Go语言编写的自定义解释器，实现了词法分析器（Lexer）、解析器（Parser）、抽象语法树（AST）和求值器（Evaluator）。

本项目灵感来源于Thorsten Ball所著的《用Go语言自制解释器》一书。在原版的基础上，添加了一些新的关键字和特性修改，以满足现代编程语言的特性。

## 现在支持的功能
- 变量赋值和引用
- 函数定义和调用
- 算术运算
- 条件语句
- 循环语句 [While]
- 数组操作
- 导入外部模块
- 注释

## 当前问题
- 错误报告缺乏行和列信息

## 未来工作
- 宏系统

## 示例用法

你可以通过一个简单的例子来学习这门新语言：

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

## REPL和文件执行支持

要在REPL模式下运行解释器或执行文件，请使用以下命令：

```bash
go run main.go -f <file_path> -t
```

`t` 是一个可选参数，用于启用执行过程的时间测量。

本项目旨在提供一个学习平台，用于构建解释器和理解编程语言设计的复杂性。欢迎贡献和反馈！