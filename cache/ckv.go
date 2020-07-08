package cache

/*
    @Time : 2020-7-8 08:05 下午
    @Author : jake
    手动实现的一个 内存存储, 尝试解决 缓存遇到的 部分问题
    * 可以添加 Redis的 部分 MaxMemory 策略
    * 热数据保持
    * 缓存污染
    * 缓存穿透
    * 持久化策略
    * 分布式怎么解决？
    * raft协议？
*/

var KeyNotFoundError = errors.New("Key not found ")

type ICache interface{
    Set(key string, value interface{}) error
    SetWithExpire(key string,value interface{},exp time.Duration) error
    Expire(key string,exp time.Duration) error
    Get(key string) (interface{}, error)
    GetAll() map[string]interface{}
    Exist(key string) bool
    Keys() []string
    Len() int
}

type Storage struct{
    defaultExpiration time.Duration
    map[string]Item // 存储实体,不存指针,防止GC扫描整个map
    lock sync.RWMutex // 表级锁

}

type Item struct{
    Obj interface{} // 存储对象
    Expiration int64 // 过期时间
    lock sync.RWMutex // 行级锁
}

func (it Item) Expired() bool {
    if item.Expiration == 0 { // 持久存储,永不过期
        return false
    }
    return time.Now().UnixNano() > item.Expiration
}



