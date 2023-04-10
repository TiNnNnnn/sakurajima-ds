package config_server


type ConfigStateMachine interface {
	Join(groups map[int][]string) error
	Leave(groupId []int) error
	Move(bucketId int, groupId int) error
	Query(num int) (Config, error)
}


