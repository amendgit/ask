---
id: first-responder
title: null
tags:
	- null
---

<!--front-->
iOS系统是如何查找第一响应者的？

<!--back-->
基础API

查找第一响应者时，有两个非常关键的API，查找第一响应者就是通过不断调用子视图的这两个API完成的。

调用方法，获取到被点击的视图，也就是第一响应者。

- (UIView *)hitTest:(CGPoint)point withEvent:(UIEvent *)event;

hitTest:withEvent:方法内部会通过调用这个方法，来判断点击区域是否在视图上，是则返回YES，不是则返回NO。

- (BOOL)pointInside:(CGPoint)point withEvent:(UIEvent *)event;

查找第一响应者

应用程序接收到事件后，将事件交给keyWindow并转发给根视图，根视图按照视图层级逐级遍历子视图，并且遍历的过程中不断判断视图范围，并最终找到第一响应者。

从keyWindow开始，向前逐级遍历子视图，不断调用UIView的hitTest:withEvent:方法，通过该方法查找在点击区域中的视图后，并继续调用返回视图的子视图的hitTest:withEvent:方法，以此类推。如果子视图不在点击区域或没有子视图，则当前视图就是第一响应者。

在hitTest:withEvent:方法中，会从上到下遍历子视图，并调用subViews的pointInside:withEvent:方法，来找到点击区域内且最上面的子视图。如果找到子视图则调用其hitTest:withEvent:方法，并继续执行这个流程，以此类推。如果子视图不在点击区域内，则忽略这个视图及其子视图，继续遍历其他视图。

可以通过重写对应的方法，控制这个遍历过程。通过重写pointInside:withEvent:方法，来做自己的判断并返回YES或NO，返回点击区域是否在视图上。通过重写hitTest:withEvent:方法，返回被点击的视图。
此方法在遍历视图时，忽略以下三种情况的视图，如果视图具有以下特征则忽略。但是视图的背景颜色是clearColor，并不在忽略范围内。

* 视图的hidden等于YES。
* 视图的alpha小于等于0.01。
* 视图的userInteractionEnabled为NO。

如果点击事件是发生在视图外，但在其子视图内部，子视图也不能接收事件并成为第一响应者。这是因为在其父视图进行hitTest:withEvent:的过程中，就会将其忽略掉。
