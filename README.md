Crust "Programming Language"
============================

Crust is a toy programming language created for the sole purpose of seeing what it would be like to create my own programming language.

The crust language itself is similar to Java bytecode, but it is very slimmed down.

You probably shouldn't use this language for any real world purposes.

### Examples

Hello world:

```
spush hello
spush  -
spush world
sadd
sadd
put
```

Counting:

```
ipush 0
ipush 1
iadd
dup
put
putln
dup
jumpl 10 2
```

This counts from 1 to 10 (inclusive) and prints each number on the screen.