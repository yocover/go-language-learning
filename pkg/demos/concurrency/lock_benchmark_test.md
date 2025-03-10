
```shell
wangzhongjie@mac concurrency % go test -bench=. -benchmem -cpu=1,2,4,8
goos: darwin
goarch: arm64
pkg: example/features/concurrency
cpu: Apple M2 Pro
BenchmarkMutex_ReadMostly               87045886                13.76 ns/op            0 B/op          0 allocs/op
BenchmarkMutex_ReadMostly-2             28968073                41.87 ns/op            0 B/op          0 allocs/op
BenchmarkMutex_ReadMostly-4             12156093                95.63 ns/op            0 B/op          0 allocs/op
BenchmarkMutex_ReadMostly-8              9381174               128.7 ns/op             0 B/op          0 allocs/op
BenchmarkRWMutex_ReadMostly             85440673                13.94 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutex_ReadMostly-2           44999929                25.91 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutex_ReadMostly-4           33925098                35.78 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutex_ReadMostly-8           40682792                37.13 ns/op            0 B/op          0 allocs/op
BenchmarkMutex_EqualReadWrite           91070070                13.11 ns/op            0 B/op          0 allocs/op
BenchmarkMutex_EqualReadWrite-2         28805444                40.65 ns/op            0 B/op          0 allocs/op
BenchmarkMutex_EqualReadWrite-4         11237392               105.7 ns/op             0 B/op          0 allocs/op
BenchmarkMutex_EqualReadWrite-8          8834019               137.1 ns/op             0 B/op          0 allocs/op
BenchmarkRWMutex_EqualReadWrite         71798964                16.94 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutex_EqualReadWrite-2       32197153                36.81 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutex_EqualReadWrite-4       22296422                53.45 ns/op            0 B/op          0 allocs/op
BenchmarkRWMutex_EqualReadWrite-8       20959249                58.46 ns/op            0 B/op          0 allocs/op
PASS
ok      example/features/concurrency    21.245s
```

让我分析一下基准测试的结果。这个测试在不同的 CPU 核心数（1,2,4,8）和不同的读写比例下比较了互斥锁（Mutex）和读写锁（RWMutex）的性能：

### 1. 读多写少场景（90% 读，10% 写）

```
# 单核性能
BenchmarkMutex_ReadMostly      87045886    13.76 ns/op
BenchmarkRWMutex_ReadMostly    85440673    13.94 ns/op

# 8核性能
BenchmarkMutex_ReadMostly-8     9381174   128.7 ns/op
BenchmarkRWMutex_ReadMostly-8   40682792   37.13 ns/op
```

分析：
1. 单核情况下，两种锁性能相近
2. 8核情况下，读写锁性能明显优于互斥锁：
   - 互斥锁：128.7 ns/op
   - 读写锁：37.13 ns/op（快约3.5倍）

### 2. 读写均衡场景（50% 读，50% 写）

```
# 单核性能
BenchmarkMutex_EqualReadWrite      91070070    13.11 ns/op
BenchmarkRWMutex_EqualReadWrite    71798964    16.94 ns/op

# 8核性能
BenchmarkMutex_EqualReadWrite-8     8834019   137.1 ns/op
BenchmarkRWMutex_EqualReadWrite-8   20959249   58.46 ns/op
```

分析：
1. 单核情况下，互斥锁略优于读写锁
2. 8核情况下，读写锁仍然优于互斥锁：
   - 互斥锁：137.1 ns/op
   - 读写锁：58.46 ns/op（快约2.3倍）

### 3. 扩展性对比

随着核心数增加的性能变化：

**读多写少场景（Mutex）**：
```
1核：13.76 ns/op
2核：41.87 ns/op （变慢 3倍）
4核：95.63 ns/op （变慢 7倍）
8核：128.7 ns/op （变慢 9.4倍）
```

**读多写少场景（RWMutex）**：
```
1核：13.94 ns/op
2核：25.91 ns/op （变慢 1.9倍）
4核：35.78 ns/op （变慢 2.6倍）
8核：37.13 ns/op （变慢 2.7倍）
```

### 结论和使用建议：

1. **单核场景**：
   - 两种锁性能相近
   - 互斥锁实现更简单，可以优先使用

2. **多核场景**：
   - 读写锁在多核下扩展性明显更好
   - 特别是在读多写少的场景下，读写锁优势更大

3. **适用场景**：
   - 读多写少场景：强烈推荐使用读写锁（RWMutex）
   - 读写均衡场景：
     - 单核下使用互斥锁（Mutex）
     - 多核下使用读写锁（RWMutex）
   - 写多读少场景：使用互斥锁（Mutex）

4. **注意事项**：
   - 读写锁实现更复杂，需要正确区分读写操作
   - 在并发量较小或单核场景下，互斥锁可能是更简单的选择
   - 如果写操作比例高，读写锁的优势会减小

这个基准测试很好地展示了为什么在实际应用中（特别是像缓存这样的读多写少场景），读写锁是更好的选择。它能在保证数据安全的同时，提供更好的并发性能。
