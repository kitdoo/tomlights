// Package storages groups the persistence layer. Each aggregate gets its
// own subpackage with an interface.go (Storage interface) and a mongo/
// implementation built on the shared *datamongo.Mongo client provided by
// internal/fx.InfrastructureModule, e.g.:
//
//	internal/storages/products/interface.go
//	internal/storages/products/mongo/storage.go
//	internal/storages/products/mongo/<unix-ts>_products_indexes.go
package storages
