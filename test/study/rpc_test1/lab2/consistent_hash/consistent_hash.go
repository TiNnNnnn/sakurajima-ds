package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/*
replicateCount: 每台服务所对应的节点数量（实际节点 + 虚拟节点）
node: 键：节点哈希值 ， 值：服务器地址
sortedNodes: 哈希环 从小到大排序后的所有节点哈希值切片
*/
type HashRing struct {
	replicateCount int
	nodes          map[uint32]string
	sortedNodes    []uint32
}

/*
构造一个哈希环
*/
func MakeHashRing(rpliCount int) *HashRing {
	newHashRing := &HashRing{
		replicateCount: rpliCount,
		nodes:          make(map[uint32]string),
		sortedNodes:    make([]uint32, 0),
	}
	return newHashRing
}

/*
 * 作用：在哈希环上添加单个服务器节点（以及虚拟节点）
 * 参数：服务器地址
 *
 * 要求：为每台服务器生成数量为 replicateCount-1 个虚拟节点，并将其与服务器的实际节点一同
 *		添加到哈希环中。
 *
 */
func (hr *HashRing) AddNode(masterNode string) {

	// 为每台服务器生成数量为 replicateCount-1 个虚拟节点
	// 并将其与服务器的实际节点一同添加到哈希环中
	for i := 0; i < hr.replicateCount; i++ {
		// 获取节点的哈希值，其中节点的字符串为 i+address
		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		// 设置该节点所对应的服务器（建立节点与服务器地址的映射）
		hr.nodes[key] = masterNode
		// 将节点的哈希值添加到哈希环中
		hr.sortedNodes = append(hr.sortedNodes, key)
	}

	// 按照值从大到小的排序函数
	sort.Slice(hr.sortedNodes, func(i, j int) bool {
		return hr.sortedNodes[i] < hr.sortedNodes[j]
	})
}

/*
 * using：添加多个服务器节点（包含虚拟节点）
 * output：服务器地址集合
 */
func (hr *HashRing) AddNodes(masterNodes []string) {
	if len(masterNodes) > 0 {
		for _, node := range masterNodes {
			// 调用 addNode 方法为每台服务器创建实际节点和虚拟节点并建立映射关系
			// 最后将创建好的节点添加到哈希环中
			hr.AddNode(node)
		}
	}
}

/*
 * using：从哈希环上移除单个服务器节点（包含虚拟节点）
 * args：服务器地址
 */
func (hr *HashRing) RemoveNode(masterNode string) {

	// 移除时需要将服务器的实际节点和虚拟节点一同移除
	for i := 0; i < hr.replicateCount; i++ {
		// 计算节点的哈希值
		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		// 移除映射关系
		delete(hr.nodes, key)
		// 从哈希环上移除实际节点和虚拟节点
		if success, index := hr.getIndexForKey(key); success {
			hr.sortedNodes = append(hr.sortedNodes[:index], hr.sortedNodes[index+1:]...)
		}
	}
}

/*
 * using：给定一个客户端地址获取应当处理其请求的服务器的地址
 * args：客户端地址
 * reply：应当处理该客户端请求的服务器的地址
 */
func (hr *HashRing) GetNode(key string) string {

	if len(hr.nodes) == 0 {
		return ""
	}

	// 获取客户端地址的哈希值
	hashKey := hr.hashKey(key)
	nodes := hr.sortedNodes

	// 当客户端地址的哈希值大于服务器上所有节点的哈希值时默认交给首个节点处理
	masterNode := hr.nodes[nodes[0]]

	for _, node := range nodes {
		// 如果客户端地址的哈希值小于当前节点的哈希值
		// 说明客户端的请求应当由该节点所对应的服务器来进行处理（逆时针）
		if hashKey < node {
			masterNode = hr.nodes[node]
			break
		}
	}

	return masterNode
}

/*
 * using：hash function,base on crc32
 * args：节点或客户端地址
 * reply：地址所对应的哈希值
 */
func (hr *HashRing) hashKey(key string) uint32 {
	scratch := []byte(key)
	return crc32.ChecksumIEEE(scratch)
}

func (hr *HashRing) getIndexForKey(key uint32) (bool, int) {
	for i, k := range hr.sortedNodes {
		if k == key {
			return true, i
		}
	}
	return false, -1
}
