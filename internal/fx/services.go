package fx

import (
	"go.uber.org/fx"

	"github.com/kitdoo/tomlights/internal/pkg/appconfig"
	"github.com/kitdoo/tomlights/internal/pkg/mailer"

	mailerservice "github.com/kitdoo/tomlights/internal/pkg/mailer/mailer"
	resendservice "github.com/kitdoo/tomlights/internal/pkg/mailer/resend"
)

// ServicesModule wires the domain layer. The email sender is provided as
// mailer.Service for any service to inject. Add aggregates here:
//
//	fx.Provide(fx.Annotate(xmongo.New, fx.As(new(x.Storage)))),
//	fx.Provide(fx.Annotate(xservice.New, fx.As(new(xsvc.Service)))),
//	fx.Provide(AsGRPCHandler(xhandler.New)),
func ServicesModule() fx.Option {
	return fx.Module("services",
		fx.Provide(fx.Annotate(newMailerService, fx.As(new(mailer.Service)))),
	)
}

// newMailerService resolves the mailer provider from cfg.App: Resend takes
// precedence over Smtp when both are set. Neither configured yields a
// Service whose Send always fails with errs.ErrMailerNotConfigured, rather
// than failing app boot.
func newMailerService(cfg *appconfig.Config) mailer.Service {
	if cfg.App == nil {
		return mailerservice.New(nil)
	}
	if resend := cfg.App.Resend; resend != nil {
		return resendservice.New(&resendservice.Config{
			APIKey: resend.APIKey.Expose(),
			From:   resend.From,
		})
	}
	if smtp := cfg.App.Smtp; smtp != nil {
		return mailerservice.New(&mailerservice.Config{
			Host:     smtp.Host,
			Port:     smtp.Port,
			Username: smtp.Username,
			Password: smtp.Password.Expose(),
			From:     smtp.From,
		})
	}
	return mailerservice.New(nil)
}
