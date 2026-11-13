package appconfig

import (
	"github.com/altessa-s/go-atlas/config"
	"github.com/altessa-s/go-atlas/core/types/redacted"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// maxPort is the highest valid TCP port number.
const maxPort = 65535

// AppConfig holds business-domain configuration that has no go-atlas infra
// counterpart.
type AppConfig struct {
	// Smtp configures the outbound mail server the email sender delivers
	// through (internal/pkg/mailer) when Resend is absent. Optional: when
	// both are absent, sending returns errs.ErrMailerNotConfigured instead
	// of failing app boot.
	Smtp *SMTPConfig `yaml:"smtp" default:"-"`

	// Resend configures delivery through the Resend HTTP API
	// (internal/pkg/mailer/resend). Takes precedence over Smtp when both
	// are set — direct SMTP from cloud/VPS IPs is often blackholed, an
	// HTTPS API call is not.
	Resend *ResendConfig `yaml:"resend" default:"-"`
}

// Validate cascades into the configured sub-sections.
func (c *AppConfig) Validate() error {
	return config.ValidateStruct(c,
		validation.Field(&c.Smtp, validation.NilOrNotEmpty),
		validation.Field(&c.Resend, validation.NilOrNotEmpty),
	)
}

// SMTPConfig configures the outbound SMTP server.
type SMTPConfig struct {
	// Host is the SMTP server hostname (e.g. "smtp.gmail.com").
	Host string `yaml:"host"`
	// Port is the SMTP server port. 587 (STARTTLS) and 465 (implicit TLS)
	// are both supported.
	Port int `yaml:"port" default:"587"`
	// Username authenticates to the SMTP server. Empty disables auth.
	Username string `yaml:"username"`
	// Password authenticates to the SMTP server alongside Username.
	Password redacted.RedactedString `yaml:"password"`
	// From is the envelope/header From address emails are sent as.
	From string `yaml:"from"`
}

// Validate requires Host, Port, and From.
func (s *SMTPConfig) Validate() error {
	return config.ValidateStruct(s,
		validation.Field(&s.Host, validation.Required),
		validation.Field(&s.Port, validation.Required, validation.Min(1), validation.Max(maxPort)),
		validation.Field(&s.From, validation.Required),
	)
}

// ResendConfig configures outbound mail delivery through the Resend HTTP API.
type ResendConfig struct {
	// APIKey authenticates to the Resend API.
	APIKey redacted.RedactedString `yaml:"apiKey"`
	// From is the header From address. Must be on a domain verified with Resend.
	From string `yaml:"from"`
}

// Validate requires APIKey and From.
func (r *ResendConfig) Validate() error {
	return config.ValidateStruct(r,
		validation.Field(&r.APIKey, validation.Required),
		validation.Field(&r.From, validation.Required),
	)
}
