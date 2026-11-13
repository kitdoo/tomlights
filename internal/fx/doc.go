// Package fx wires the service using Uber Fx dependency injection.
//
// It is organized into three modules:
//
//   - [InfrastructureModule] — MongoDB, TLS, metrics, tracing, health coordinator
//   - [ServicesModule] — domain storages, services, and gRPC handlers
//   - [TransportsModule] — gRPC and HTTP servers with handler registration
//
// The infrastructure constructors return `nil, nil` when their config section
// is absent, so optional subsystems can be left out of `configs/*.yaml`.
package fx
