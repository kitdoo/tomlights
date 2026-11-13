// Package commands defines the root Cobra CLI tree for the tomlights binary.
//
// The tree is assembled in [New] and currently looks like:
//
//	tomlights       (root)
//	├── server            server lifecycle
//	│   └── run           start the gRPC + HTTP server
//	└── version           print build / version info
//
// Call [Run] to build the tree, set CLI args, and execute.
package commands
