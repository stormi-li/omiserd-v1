package omiserd

import "time"

const namespace_separator = ":"
const config_expire_time = 2 * time.Second

const prefix_Config = "stormi:config:"
const prefix_Server = "stormi:server:"
const prefix_Web = "stormi:web:"

type NodeType string

const Server NodeType = "Server"
const Config NodeType = "Config"
const Web NodeType = "Web"

const Command_update_weight = "update_weight"
