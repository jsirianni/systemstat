package postgres

import (
	"github.com/pkg/errors"
)

func (p Postgres) Validate() error {
	errs := errors.New("Failed to validate postgres configuration")
	valid := true

	if err := p.validateHost(); err != nil {
		errs = errors.Wrap(errs, err.Error())
		valid = false
	}

	if err := p.validPort(); err != nil {
		errs = errors.Wrap(errs, err.Error())
		valid = false
	}

	if err := p.validUser(); err != nil {
		errs = errors.Wrap(errs, err.Error())
		valid = false
	}

	if err := p.validPassword(); err != nil {
		errs = errors.Wrap(errs, err.Error())
		valid = false
	}

	if err := p.validDBName(); err != nil {
		errs = errors.Wrap(errs, err.Error())
		valid = false
	}

	if !valid {
		return errs
	}
	return nil
}

func (p Postgres) validateHost() error {
	if p.host == "" {
		return errors.New(envHost + " is not set")
	}
	return nil
}

func (p Postgres) validPort() error {
	if p.port == 0 {
		return errors.New(envPort + " is not set")
	}
	return nil
}

func (p Postgres) validUser() error {
	if p.user == "" {
		return errors.New(envUser + " is not set")
	}
	return nil
}

func (p Postgres) validPassword() error {
	if p.pass == "" {
		return errors.New(envPass + " is not set")
	}
	return nil
}

func (p Postgres) validDBName() error {
	if p.dbname == "" {
		return errors.New(envDBName + " is not set")
	}
	return nil
}
