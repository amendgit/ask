---
id: ios-event-transfer
title: null
draft: true
tags:
    - null
---

<!--front-->
描述iOS中事件传递。

<!--back-->
传递过程

1.UIApplication接收到事件，将事件传递给keyWindow。
2.keyWindow遍历subviews的hitTest:withEvent:方法，找到点击区域内合适的视图来处理事件。
3.UIView的子视图也会遍历其subviews的hitTest:withEvent:方法，以此类推。
4.直到找到点击区域内，且处于最上方的视图，将视图逐步返回给UIApplication。
5.在查找第一响应者的过程中，已经形成了一个响应者链。
6.应用程序会先调用第一响应者处理事件。
7.如果第一响应者不能处理事件，则调用其nextResponder方法，一直找响应者链中能处理该事件的对象。
8.最后到UIApplication后仍然没有能处理该事件的对象，则该事件被废弃。

模拟代码
```object-c
- (UIView *)hitTest:(CGPoint)point withEvent:(UIEvent *)event {
    if (self.alpha <= 0.01 || self.userInteractionEnabled == NO || self.hidden) {
        return nil;
    }
    
    BOOL inside = [self pointInside:point withEvent:event];
    if (inside) {
        NSArray *subViews = self.subviews;
        // 对子视图从上向下找
        for (NSInteger i = subViews.count - 1; i >= 0; i--) {
            UIView *subView = subViews[i];
            CGPoint insidePoint = [self convertPoint:point toView:subView];
            UIView *hitView = [subView hitTest:insidePoint withEvent:event];
            if (hitView) {
                return hitView;
            }
        }
        return self;
    }
    return nil;
}
```

示例


如上图所示，响应者链如下：

1. 如果点击UITextField后其会成为第一响应者。
2. 如果textField未处理事件，则会将事件传递给下一级响应者链，也就是其父视图。
3. 父视图未处理事件则继续向下传递，也就是UIViewController的View。
4. 如果控制器的View未处理事件，则会交给控制器处理。
5. 控制器未处理则会交给UIWindow。
6. 然后会交给UIApplication。
7. 最后交给UIApplicationDelegate，如果其未处理则丢弃事件。

事件通过UITouch进行传递，在事件到来时，第一响应者会分配对应的UITouch，UITouch会一直跟随着第一响应者，并且根据当前事件的变化UITouch也会变化，当事件结束后则UITouch被释放。

UIViewController没有hitTest:withEvent:方法，所以控制器不参与查找响应视图的过程。但是控制器在响应者链中，如果控制器的View不处理事件，会交给控制器来处理。控制器不处理的话，再交给View的下一级响应者处理。

注意

1. 在执行hitTest:withEvent:方法时，如果该视图是hidden等于NO的那三种被忽略的情况，则改视图返回nil。
2. 如果当前视图在响应者链中，但其没有处理事件，则不考虑其兄弟视图，即使其兄弟视图和其都在点击范围内。
3. UIImageView的userInteractionEnabled默认为NO，如果想要UIImageView响应交互事件，将属性设置为YES即可响应事件。