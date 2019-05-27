---
id: goroutine
title: 什么是goroutine?
draft: true
tags:
    - null
---

<!--front-->
简要描述什么是goroutine?
goroutine和进程线程有什么区别？

<!--back-->

Goroutine是一个简单的模型：它是一个函数，与其他Goroutines并发执行且共享相同地址空间。
* 通常用法是根据需要创建尽可 能的Groutines，成百上千甚至上万的
* 以多路复用的形式运行于操作系统为应用程序分配的少数几个线程上
* 创建一个Goroutine并不需要太多内存，只需要8K的栈空间，它们根据需要在堆上分配和释放内存以实现自身的增长
* 当另一个Goroutine被调度时，只需要保存/恢复三个寄存器，分别是PC、SP和DX。
* Goroutine是廉价的，更关键地是，如果它们在网络输入操作、Sleep操作、Channel操作或 sync包的原语操作上阻塞了，也不会导致承载其多路复用的线程阻塞。

参考
* [Goroutine是如何工作的](https://tonybai.com/2014/11/15/how-goroutines-work/)