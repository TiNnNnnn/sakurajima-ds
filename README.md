分布式对象存储

1. 上传流程

2. 读取流程
 - 根据bucketName+object_name从metaserver中查找到对应的objectId
 - 从metaserver的Object表中依靠objectId查找到对应的blocklist
 - 根据用户请求的读取区间，计算出要读取的block个数
 - 以blockId为key,作md5sum(blockId)%2233计算出sharedId
 - 根据shardeId查找路由服务器中的shard表,得到对应的存储集群addr
 - 根据addr,网关向sharedServer发起读写请求
 - 网关对block进行截断，拼装，返回给用户