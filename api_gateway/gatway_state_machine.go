package api_gateway

import "sakurajima-ds/storage_engine"

type ApiGatewayOp interface{

}

type ApiGatewayStateMachine struct{
	engine storage_engine.KvStorage
	
}