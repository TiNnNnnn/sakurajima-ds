package main

import (
	"fmt"
	"hash/crc32"
	ch "sakurajima-ds/test/study/rpc_test1/lab2/consistent_hash"
	"sort"
	"strconv"
)

func testAddNode() {
	hr := ch.MakeHashRing(4)
	hr.AddNode("127.0.0.1:8080")

	addr := hr.GetNode("23121321")

	if addr == "127.0.0.1:8080" {
		fmt.Println("OK")
	}
}

func testAddNodes() {
	hr := ch.MakeHashRing(4)
	testhr := MakeHashRing(4)
	masterNodes := []string{"127.0.0.1:8080", "35.57.43.10:10010", "49.235.90.234:9090"}
	hr.AddNodes(masterNodes)
	testhr.AddNodes(masterNodes)

	key := "232434#$"
	if hr.GetNode(key) == testhr.GetNode(key) {
		fmt.Println("OK")
	}
}

func testRemoveNode() {
	hr := ch.MakeHashRing(4)
	masterNodes := []string{"127.0.0.1:8080", "35.57.43.10:10010", "49.235.90.234:9090"}
	hr.AddNodes(masterNodes)

	hr.RemoveNode("49.235.90.234:9090")

	addr := hr.GetNode("23121321")
	addr2 := hr.GetNode("test43242340")

	if addr == "35.57.43.10:10010" && addr2 == "127.0.0.1:8080" {
		fmt.Println("OK")
	}
}

func main() {
	testAddNode()
	testAddNodes()
	testRemoveNode()
}


type HashRing struct {
	replicateCount int
	nodes          map[uint32]string
	sortedNodes    []uint32
}

func MakeHashRing(rpliCount int) *HashRing {
	newHashRing := &HashRing{
		replicateCount: rpliCount,
		nodes:          make(map[uint32]string),
		sortedNodes:    make([]uint32, 0),
	}
	return newHashRing
}

func (hr *HashRing) AddNode(masterNode string) {

	for i := 0; i < hr.replicateCount; i++ {

		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		hr.nodes[key] = masterNode
		hr.sortedNodes = append(hr.sortedNodes, key)
	}
	sort.Slice(hr.sortedNodes, func(i, j int) bool {
		return hr.sortedNodes[i] < hr.sortedNodes[j]
	})
}

func (hr *HashRing) AddNodes(masterNodes []string) {
	if len(masterNodes) > 0 {
		for _, node := range masterNodes {

			hr.AddNode(node)
		}
	}
}

func (hr *HashRing) RemoveNode(masterNode string) {

	for i := 0; i < hr.replicateCount; i++ {

		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		delete(hr.nodes, key)
		if success, index := hr.getIndexForKey(key); success {
			hr.sortedNodes = append(hr.sortedNodes[:index], hr.sortedNodes[index+1:]...)
		}
	}
}

func (hr *HashRing) GetNode(key string) string {

	if len(hr.nodes) == 0 {
		return ""
	}

	hashKey := hr.hashKey(key)
	nodes := hr.sortedNodes

	masterNode := hr.nodes[nodes[0]]

	for _, node := range nodes {
		if hashKey < node {
			masterNode = hr.nodes[node]
			break
		}
	}

	return masterNode
}

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
