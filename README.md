cbuild
======

## What's cbuild

`cbuild` is a very easy builder for building c/c++ program.

## Getting started

First create a directory named `demo`(or any name you like). Then create a source file named `main.cpp`(or any name which has suffix `.cpp`,`.cxx`,`.hpp`,`.hxx`,`.cc` or `.c`) in directory `demo`.

```cpp
// main.cpp
#include <stdio.h>

int main() {
	printf("hello, cbuild!\n");
}
```

Finally, execute `cbuild run` in directory `demo`, it will build the c/c++ program and then execute it.

```shell
cbuild run
```

You also can execute `cbuild` instead of `cbuild run`, it will build the c/c++ program only.

See more complex example [testdata](https://github.com/mkideal/cbuild/tree/master/testdata)

Execute `cbuild -h` to display help information.

