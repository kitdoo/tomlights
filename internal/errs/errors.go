// Package errs holds sentinel errors shared across services, storages, and
// transports. Compare with errors.Is, never string equality.
package errs

import "errors"

var (
	// ErrInvalidArgument is returned when a request argument fails domain
	// validation.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrNotImplemented marks a method whose body is still a skeleton.
	ErrNotImplemented = errors.New("not implemented")

	// ErrMailerNotConfigured is returned by mailer.Service.Send when
	// neither AppConfig.Resend nor AppConfig.Smtp is set.
	ErrMailerNotConfigured = errors.New("mailer not configured")
)
