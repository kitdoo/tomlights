package fx

import (
	"context"
	"errors"
	"log/slog"

	"go.uber.org/fx"

	"github.com/altessa-s/go-atlas/config"
	"github.com/altessa-s/go-atlas/observability/health"
	"github.com/altessa-s/go-atlas/observability/metrics"
	"github.com/altessa-s/go-atlas/observability/tracing"

	"github.com/kitdoo/tomlights/internal/pkg/appconfig"

	coreerrs "github.com/altessa-s/go-atlas/core/errors"
	datamongo "github.com/altessa-s/go-atlas/data/mongo"
	mongofactory "github.com/altessa-s/go-atlas/infrastructure/mongo/factory"
	metricsfactory "github.com/altessa-s/go-atlas/observability/metrics/factory"
	slogx "github.com/altessa-s/go-atlas/observability/slog"
	tracingfactory "github.com/altessa-s/go-atlas/observability/tracing/factory"
	tlsfactory "github.com/altessa-s/go-atlas/security/tlsutils/factory"
	tlsproviders "github.com/altessa-s/go-atlas/security/tlsutils/providers"
)

// mongoHealthService is the name the MongoDB checker is registered under
// in the health coordinator (visible in /internal/healthz/details).
const mongoHealthService = "mongodb"

// InfrastructureModule provides core infrastructure components: MongoDB,
// TLS, metrics, tracing, and health.
//
// Constructors that depend on an optional configuration section return
// (nil, nil) when the section is absent — downstream consumers MUST nil-check
// before using these values.
func InfrastructureModule() fx.Option {
	return fx.Module("infrastructure",
		fx.Provide(newMongo),
		fx.Provide(newTLSProviders),
		fx.Provide(newMetricsCollector),
		fx.Provide(newTracer),
		fx.Provide(newHealthCoordinator),
		fx.Invoke(registerMongoHealth),
	)
}

func newMongo(cfg *appconfig.Config, collector metrics.Collector, lc fx.Lifecycle) (*datamongo.Mongo, error) {
	initCtx, cancel := initContext()
	defer cancel()

	m, err := mongofactory.New(cfg.MongoDB).
		UseLogger(slog.Default().With(slogx.Module("mongodb"))).
		UseCollector(collector).
		Build(initCtx)
	if err != nil {
		return nil, coreerrs.WrapOperation(err, "create MongoDB client")
	}

	if err = m.Connect(initCtx); err != nil {
		return nil, coreerrs.WrapOperation(err, "connect to MongoDB")
	}

	lc.Append(fx.StopHook(m.Close))
	return m, nil
}

// registerMongoHealth makes /internal/readyz (and the gRPC health service)
// report NotServing while MongoDB is unreachable.
func registerMongoHealth(hc *health.Coordinator, m *datamongo.Mongo) {
	hc.RegisterService(mongoHealthService, health.Func(func(ctx context.Context) health.ServingStatus {
		if err := m.Client().Ping(ctx, nil); err != nil {
			return health.StatusNotServing
		}
		return health.StatusServing
	}))
}

func newTLSProviders(cfg *appconfig.Config, lc fx.Lifecycle) (*tlsproviders.Providers, error) {
	if !cfg.IsTlsProviderConfigured() {
		return nil, nil //nolint:nilnil // TLS provider is optional
	}

	s, err := tlsfactory.New(cfg.TlsProvider).
		UseLogger(slog.Default().With(slogx.Module("tls"))).
		Build()
	if err != nil {
		return nil, coreerrs.WrapOperation(err, "create tls providers")
	}

	lc.Append(fx.StopHook(func(ctx context.Context) {
		s.Close(ctx, nil)
	}))
	return s, nil
}

func newMetricsCollector(cfg *appconfig.Config, lc fx.Lifecycle) (metrics.Collector, error) {
	var metricsCfg *config.Metrics
	if cfg.Observability != nil {
		metricsCfg = cfg.Observability.Metrics
	}

	collector, err := metricsfactory.New(metricsCfg).Build()
	if err != nil {
		return nil, coreerrs.WrapOperation(err, "create metrics collector")
	}

	lc.Append(fx.StopHook(collector.Shutdown))
	return collector, nil
}

func newTracer(cfg *appconfig.Config, lc fx.Lifecycle) (tracing.Tracer, error) {
	var tracingCfg *config.Tracing
	if cfg.Observability != nil {
		tracingCfg = cfg.Observability.Tracing
	}

	initCtx, cancel := initContext()
	defer cancel()

	tracer, err := tracingfactory.New(tracingCfg).Build(initCtx)
	if err != nil {
		return nil, coreerrs.WrapOperation(err, "create tracer")
	}

	lc.Append(fx.StopHook(tracer.Shutdown))
	return tracer, nil
}

func newHealthCoordinator(collector metrics.Collector) *health.Coordinator {
	return health.New(health.WithCollector(collector))
}

var errServiceNotReady = errors.New("service is not ready")

// healthAdapter adapts [health.Coordinator] to the gRPC server's health.Health interface.
type healthAdapter struct {
	coordinator *health.Coordinator
}

// Health reports the current readiness of the service.
func (a *healthAdapter) Health(ctx context.Context) error {
	if a.coordinator.CheckHealth(ctx) != health.StatusServing {
		return errServiceNotReady
	}
	return nil
}
