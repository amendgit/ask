package main

import (
	"testing"
)

func TestCard(t *testing.T) {
	cardDAO := NewCardDAO()
	card := cardDAO.ParseString(`---
id: acid
title: DBMS中的ACID指的是什么？
tags:
    - database
    - interview
---
	
<!--front-->
DBMS中的ACID指的是什么？

<!--back-->
ACID，是指在可靠数据库管理系统（DBMS）中，事务(transaction)所应该具有的四个特性：
* 原子性（Atomicity）: 原子性是指事务是一个不可再分割的工作单位，事务中的操作要么都发生，要么都不发生。
* 一致性（Consistency）：一致性是指在事务开始之前和事务结束以后，数据库的完整性约束没有被破坏。
* 隔离性（Isolation）：多个事务并发访问时，事务之间是隔离的，一个事务不应该影响其它事务运行效果。
* 持久性（Durability）：持久性，意味着在事务完成以后，该事务所对数据库所作的更改便持久的保存在数据库之中，并不会被回滚。
这是可靠数据库所应具备的几个特性.	
`)
	got := card.Metadata.String()
	expect := `tags: [database,interview]
title: DBMS中的ACID指的是什么？
`
	if got != expect {
		t.Fatalf("got %v expect %v", got, expect)
	}
}
