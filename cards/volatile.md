---
id   : volatile
title: volatile
---

<!--front-->

如何理解volitale关键字？

<!--back-->

volatile是一个类型修饰符（type specifier）.volatile的作用是作为指令关键字，确保本条指令不会因编译器的优化而省略，且要求每次直接读值。 volatile的变量是说这变量可能会被意想不到地改变，这样，编译器就不会去假设这个变量的值了。