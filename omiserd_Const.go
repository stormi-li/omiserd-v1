package omiserd

import "time"

const Namespace_separator = ":"
const Config_expire_time = 2 * time.Second

const Prefix_Config = "stormi:config:"
const Prefix_Server = "stormi:server:"
const Prefix_Web = "stormi:web:"

type NodeType string

const Server NodeType = "Server"
const Config NodeType = "Config"
const Web NodeType = "Web"

const Command_update_weight = "update_weight"
