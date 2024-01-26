# 一、为什么选择缓存
提升性能，快速响应。通常我们对于持久化数据都会放在数据库中，一般来说，我们都是直接访问数据库的，数据的性能瓶颈在于

- 网络连接
- 磁盘io，取决于sql 性能与数据量等等。 一般我们针对我们的应用程序会增加一层缓存，在选择缓存时，我们也会有很多选择、包括redis、memcache、这些都是作为进程独立部署的，可以用于分布式缓存。
  如果从业务量出发，同时成本低，我们可以采用本地缓存的方式，本地缓存具有的优点
- 访问快，不需要网络连接
- 易于实现 缺点
- 不能分布式存储，也就是每台服务器都会利用一段内存数据。 缓存通常用于存储我们的热点数据，通常缓存数据具有的特点：访问频率高、不易修改、数据量有限。

# 二、如何设计一个本地缓存
## 1.freecatch 
1.1内存预分配 
1.2 并发安全
1.3 采用分段的方式，将预分配的内存，分到256个segment,每个segment 会利用ringBuf数据结构存储，此种方式的弊端，预分配会浪费内存，无法扩容，无法存在较大的数据，这是该框架的最大弊端。
```go
cacheSize := 100 * 1024 * 1024 //100MB
cacheSize := 100 * 1024 * 1024
cache := freecache.NewCache(cacheSize)
debug.SetGCPercent(20)
key := []byte("abc")
val := []byte("def")
expire := 60 // expire in 60 seconds
cache.Set(key, val, expire)
got, err := cache.Get(key)
if err != nil {
fmt.Println(err)
} else {
fmt.Printf("%s\n", got)
}
affected := cache.Del(key)
fmt.Println("deleted key ", affected)
fmt.Println("entry count ", cache.EntryCount())
```
## 2、groupcatch
- 2.1需要和业务管理，独立进程
- 2.2对等模型，分布式缓存
- 2.3通常用来代替memcache
使用样例，较为复杂。由于groupcatch 是一种对等模型，因此groupcatch 在分布式存储上有优势。
  
3、bigcache
快速、并发、逐出内存缓存，用于在不影响性能的情况下保留大量条目。BigCache将条目保留在堆上，但省略它们的GC。为了实现这一点，需要对字节片进行操作，因此在大多数用例中都需要缓存前面的条目（反）序列化。
BigCache依赖于Go 1.5版本（issue-9477）中提供的优化。此优化声明，如果使用键和值中没有指针的映射，则GC将省略其内容。因此，BigCache使用map[uint64]
uint32，其中键是散列的，值是条目的偏移量。条目保存在字节片中，以再次省略GC。字节片大小可以增长到千兆字节，而不会影响性能，因为GC只能看到指向它的单个指针。BigCache不处理冲突。当插入新项并且其哈希与以前存储的项冲突时，新项将覆盖以前存储的值。

# 三、本地缓存的简单实现
基于不同的业务场景，我们可以选择设计一个符合自己业务场景的缓存。 对于我需要的而言，我需要是 
1、能够定时从数据库中自动加载数据。 
2、支持并发安全 
3、支持过期淘汰 
4、便于扩展、较为通用

# 四、利用bigcache实现上述场景

# 五、数据库表容量
```sql
select table_schema as '数据库', table_name as '表名', table_rows as '记录数', truncate(data_length / 1024 / 1024, 2) as '数据容量(MB)', truncate(index_length / 1024 / 1024, 2) as '索引容量(MB)'
from information_schema.tables
where table_schema = 'saas'
  and table_name = 'saas_tenant'
order by data_length desc, index_length desc;
```
```sql
select table_schema as '数据库', table_name as '表名', table_rows as '记录数', truncate(data_length / 1024 / 1024, 2) as '数据容量(MB)', truncate(index_length / 1024 / 1024, 2) as '索引容量(MB)'
from information_schema.tables
where table_schema = 'sentry'
  and table_name = 'sentry_rule_engine_app'
order by data_length desc, index_length desc;
```

#六、业务实现
```go
service.InitOrganizationCache(ctx,
service.ICacheConfig{
LoadFunc: model.GetMap,
CacheName: "公司",
LoadInterval: 20 * time.Minute,
})
service.InitAppCache(ctx, service.ICacheConfig{
LoadFunc: model.GetAppMap,
CacheName: "应用",
LoadInterval: 20 * time.Minute,
})
service.ICacheReLoadFunc(ctx, service.OrganizationCache)
service.ICacheReLoadFunc(ctx, service.AppCache)
```