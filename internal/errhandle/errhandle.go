package errhandle

import (
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/internal/slack"
)

func Handle(ctx context.Ctx, aerr errapi.Error, err error, fatal bool) {
	if aerr != nil {
		ctx = context.WithValue(ctx, context.KeyErrorAPI, aerr.Code())
	}

	if err != nil {
		slack.Client.MessageError(ctx, err)
		ctx = context.WithValue(ctx, context.KeyError, err)
	}

	if fatal {
		logger.Logger.Emerf(ctx, `err="emergency exit"`)
		logger.Logger.Fatalf(`err="%v"`, err)
	} else {
		logger.Logger.Errorf(ctx, "")
	}
}
