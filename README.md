# bfg
**A fast Brainfuck interpreter written in the Go programming language.**

---

## Overview

Brainfuck is an esoteric programming language noted for its extreme minimalism. The language consists of only eight simple commands and an instruction pointer. It is designed to challenge and amuse programmers, and was not made to be suitable for practical use. It was created in 1993 by Urban Müller. [See wikipedia page.](https://en.wikipedia.org/wiki/Brainfuck)

### Usage

Brainfuck programs can be run by specifing the `-f` switch.

```
bfg -f program.bf
```

### Linux

On linux a platform specific compiler is used by default. This means brainfuck code is executed up to five times faster than other platforms. To revert to use the cross-platform interpreter, pass the `-i` switch.

```
bfg -i -f program.bf
```
