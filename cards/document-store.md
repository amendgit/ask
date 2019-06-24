---
id: document-store
title: null
tags:
    - null
---

<!--front-->
文档类型存储

<!--back-->
抽象模型：将文档作为值的键-值存储

文档类型存储以文档（XML、JSON、二进制文件等）为中心，文档存储了指定对象的全部信息。文档存储根据文档自身的内部结构提供 API 或查询语句来实现查询。请注意，许多键-值存储数据库有用值存储元数据的特性，这也模糊了这两种存储类型的界限。

基于底层实现，文档可以根据集合、标签、元数据或者文件夹组织。尽管不同文档可以被组织在一起或者分成一组，但相互之间可能具有完全不同的字段。

MongoDB 和 CouchDB 等一些文档类型存储还提供了类似 SQL 语言的查询语句来实现复杂查询。DynamoDB 同时支持键-值存储和文档类型存储。

文档类型存储具备高度的灵活性，常用于处理偶尔变化的数据。

来源及延伸阅读：文档类型存储
* 面向文档的数据库：https://en.wikipedia.org/wiki/Document-oriented_database
* MongoDB 架构：https://www.mongodb.com/mongodb-architecture
* CouchDB 架构：https://blog.couchdb.org/2016/08/01/couchdb-2-0-architecture/
* Elasticsearch 架构：https://www.elastic.co/blog/found-elasticsearch-from-the-bottom-up

