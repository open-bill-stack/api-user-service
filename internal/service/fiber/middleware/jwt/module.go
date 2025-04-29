package jwt

import "go.uber.org/fx"

type Params struct {
	fx.In
}

type Result struct {
	fx.Out

	Middleware *CheckJWTMiddleware
}

func NewService(p Params) (Result, error) {
	checkJWTMiddleware := NewCheckJWTMiddleware()

	return Result{
		Middleware: checkJWTMiddleware,
	}, nil
}
