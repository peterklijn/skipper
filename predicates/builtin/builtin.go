package builtin

import (
	"github.com/zalando/skipper/predicates"
	pauth "github.com/zalando/skipper/predicates/auth"
	"github.com/zalando/skipper/predicates/cookie"
	"github.com/zalando/skipper/predicates/cron"
	"github.com/zalando/skipper/predicates/interval"
	"github.com/zalando/skipper/predicates/methods"
	"github.com/zalando/skipper/predicates/primitive"
	"github.com/zalando/skipper/predicates/query"
	"github.com/zalando/skipper/predicates/source"
	"github.com/zalando/skipper/predicates/tee"
	"github.com/zalando/skipper/predicates/traffic"
)

func MakeRegistry() predicates.Registry {
	r := make(predicates.Registry)
	for _, s := range []predicates.PredicateSpec{
		source.New(),
		source.NewFromLast(),
		source.NewClientIP(),
		interval.NewBetween(),
		interval.NewBefore(),
		interval.NewAfter(),
		cron.New(),
		cookie.New(),
		query.New(),
		traffic.New(),
		primitive.NewTrue(),
		primitive.NewFalse(),
		pauth.NewJWTPayloadAllKV(),
		pauth.NewJWTPayloadAnyKV(),
		pauth.NewJWTPayloadAllKVRegexp(),
		pauth.NewJWTPayloadAnyKVRegexp(),
		methods.New(),
		tee.New(),
	} {
		r.Register(s)
	}
	return r
}
