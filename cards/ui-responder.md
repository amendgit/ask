---
id   : ui-responder
title: ui-responder
tags: 
    - ios
---

<!--front-->

如何理解UIResponder(UI响应者接口)？

<!--back-->

UIResponder是用于响应和处理事件的抽象接口。

响应者对象（即UIResponder的实例）们组成了一个UIKit应用的事件处理的框架。许多关键的对象都是响应者，例如：UIApplcation，UIViewController以及所有的UIView对象（包括UIWindow）。当事件发生的时候，UIKit将事件分发给app的响应者对象们去处理。

事件有多种：触摸、收拾、远程控制、按压事件。如果要处理某个特定类型的事件，响应者(responder)必须要重写相应的方法。比如，要处理触摸事件，响应者相应的实现touchesBegan:withEvent:, touchesMoved:withEvent:, touchesEnded:withEvent:, 和 touchesCancelled:withEvent: 方法。在触摸的事件中，响应者使用UIKit提供的事件信息，去跟踪触摸的变化，并恰当的更新app的交互。

除了处理事件之外，UIKit的响应者们还负责将未处理的事件转发给app中的其他部分。如果一个给定的响应者不能处理某个事件，它会将事件转发给在响应者链中的下一个响应者。UIKit是动态的管理响应者的。根据预定义的规则判断下一个应该负责接收事件的对象。例如，一个view会将事件转发给它的父view，而根view则会将事件转发给他的ViewController。

响应者们还可以通过inputView来接收自定义的输入，一个明显的例子就是系统键盘。当用户点击了屏幕上的UITextView或者UITextField的时候，该view会成为第一响应者，并显示他的inputView，即系统键盘。类似的，你也可以创建自定义的input views，当其他的响应者被激活的时候。将一个自定义的input view关联到一个responder，可以通过将view赋值给responder的inputView属性来实现。

更多关于响应者和响应者链的信息，[Event Handling Guide for UIKit Apps](https://developer.apple.com/documentation/uikit?language=objc)