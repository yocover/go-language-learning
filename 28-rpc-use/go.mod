module example/userpc

go 1.24.0

require (
	example/rpc/client v0.0.0
)

replace (
	example/rpc/client => ../rpc-client
)
