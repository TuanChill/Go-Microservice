package app

import (
	"context"
	"errors"
)

func (a *App) Shutdown(_ context.Context) error {
	var errs []error
	if a.MessageQueue != nil {
		errs = append(errs, a.MessageQueue.Close())
	}
	if a.Cache != nil {
		errs = append(errs, a.Cache.Close())
	}
	if a.DB != nil {
		errs = append(errs, a.DB.Close())
	}
	return errors.Join(errs...)
}
