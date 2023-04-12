package meta_server

import "sakurajima-ds/tinnraftpb"

//Name表: [BucketName Key objectId] 
type Name struct{
	BucketName string  
	Key string 
	ObjectId int64 
}

//Object表:   [ObjectId blocks]
type Object struct{
	ObjectId int64
	Blocks *tinnraftpb.DataBlocks
}


