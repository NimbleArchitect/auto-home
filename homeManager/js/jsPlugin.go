package js

import "net/rpc"

type jsPlugin struct {
	client *rpc.Client
}
