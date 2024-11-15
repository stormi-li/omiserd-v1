package omiserd

import "time"

const namespace_separator = ":"
const config_expire_time = 2 * time.Second

const Prefix_Config = "stormi:config:"
const Prefix_Server = "stormi:server:"
const Prefix_Web = "stormi:web:"

type NodeType string

var Server NodeType = "Server"
var Config NodeType = "Config"
var Web NodeType = "Web"

const Command_update_weight = "update_weight"
