---
id: designated-init
title: null
tags:
    - null
---

<!--front-->
什么是指定初始化器（Designated Initializer）?

<!--back-->
Designated Initializer通过向父类发送初始化消息类保证对象的完全初始化，实现细节在开发者继承该类时显得非常重要。

* Designated Initializer必须调用（通过super）父类的Designated Initializer，其中NSObject对应的是[super init]
* 任何便利初始化器（Convenience Initializer），必须调用类的另外一个Initializer，最终会调用Designated Initializer
* 具有指定初始值设定项的类，必须实现父类的所有指定初始值的设定项。

关键字，NS_DESIGINATED_INITIALIZER。

* 如果子类指定了新的Designated Initializer，那么在这个初始化器内部必须调用父类的Designated Initializer。并且需要重写父类的Designated Initializer，将其指向子类新的Designated Initializer。
* 如果定义NS_DESIGNATED_INITIALIZER，大多是不想让调用者调用父类的初始化函数，只希望通过该类指定的初始化进行初始化，这时候就可以用NS_UNAVAILABLE宏。
* 避免使用new，如果使用new来创建对象的话，即使init被声明为NS_UNAVAILABLE，也不会收到编译器的警告和错误提示了。
