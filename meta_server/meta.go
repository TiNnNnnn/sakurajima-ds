package meta_server

type MetaData struct {
	UserId     int64
	BucketName string
	ObjectName   string
	Version    int64
	Size       int64
	Hash       string
	//BucketId   int64
}
