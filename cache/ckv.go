package cache

/*
    @Time : 2020-7-8 08:05 下午
    @Author : jake
    手动实现的一个 内存存储, 尝试解决 缓存遇到的 部分问题
    v1.0.0
    * 可以添加 Redis的 部分 MaxMemory 策略, 如何快速定位到即将过期的 key ? 栈 ?
    * 热数据保持 每次取数据, 给该数据延长过期时间
    * 缓存污染 增加 缓存规则 , 读取频率 达到一定级别才进行缓存
    * 缓存穿透 bloom filter记录未命中的key , 每次有新key 清除bloom filter 记录
    * 持久化策略 异步对key进行持久化存储 / 序列化反序列化协议? json / protobuf / 自制二进制序列化 / 反序列化协议

    todo:
    * 分布式怎么解决？  raft协议？
    * 横向无线扩容
    * master-slave 读写分离
*/

var KeyNotFoundError = errors.New("Key not found ")

type ICache interface{
    Set(key string, value interface{}) error
    SetWithExpire(key string,value interface{},exp time.Duration) error
    Expire(key string,exp time.Duration) error
    Get(key string) (interface{}, error)
    GetAll() map[string]interface{}
    Del(key string) (bool,error)
    Exist(key string) bool
    Keys() []string
    Len() int
}

type Storage struct{
    defaultExpiration time.Duration
    map[string]Item // 存储实体,不存指针,防止GC扫描整个map
    l *sync.RWMutex // 表级锁

}

type Item struct{
    obj interface{} // 存储对象
    expiration int64 // 过期时间
    l *sync.RWMutex // 行级锁
}

func (it Item) Expired() bool {
    if item.expiration == 0 { // 持久存储,永不过期
        return false
    }
    return time.Now().UnixNano() > item.expiration
}



