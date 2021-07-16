package predicates

import (
	"github.com/zalando/skipper/eskip"
	"github.com/zalando/skipper/helpers"
	"github.com/zalando/skipper/predicates/builtin"
	"github.com/zalando/skipper/predicates/weight"
	"net/http"
	"regexp"
	"testing"
	"time"
)

func TestArgumentConversion(t *testing.T) {
	methodPredicate := Methods(http.MethodGet, http.MethodPost)
	if len(methodPredicate.Args) != 2 {
		t.Errorf("expected 2 arguments in Methods predicate, got %d", len(methodPredicate.Args))
	}

	kvPredicate := JWTPayloadAllKV(helpers.NewKVPair("k1", "v1"), helpers.NewKVPair("k2", "v2"))
	if len(kvPredicate.Args) != 4 {
		t.Errorf("expected 4 arguments in JWTPayloadAllKV predicate, got %d", len(kvPredicate.Args))
	}

	r := regexp.MustCompile(`/\d+/`)
	kvRegexPredicate := JWTPayloadAllKVRegexp(helpers.NewKVRegexPair("k1", r), helpers.NewKVRegexPair("k2", r))
	if len(kvRegexPredicate.Args) != 4 {
		t.Errorf("expected 4 arguments in JWTPayloadAllKV predicate, got %d", len(kvRegexPredicate.Args))
	}
}

func TestPredicateCreation(t *testing.T) {
	regex := regexp.MustCompile("/[a-z]+/")
	kvPair := helpers.NewKVPair("iss", "https://accounts.google.com")
	kvRegexPair := helpers.NewKVRegexPair("iss", regex)
	registry := builtin.MakeRegistry()

	t.Run("Path()", testPredicatesWithoutSpec(Path("/skipper")))
	t.Run("PathSubtree()", testPredicatesWithoutSpec(PathSubtree("/skipper")))
	t.Run("PathRegexp()", testPredicatesWithoutSpec(PathRegexp(regex)))
	t.Run("Host()", testPredicatesWithoutSpec(Host(regex)))
	t.Run("Weight()", testPredicatesWithoutSpec(Weight(42)))
	t.Run("True()", testWithSpecFn(registry, True()))
	t.Run("False()", testWithSpecFn(registry, False()))
	t.Run("Method()", testPredicatesWithoutSpec(Method("test blalalalaal")))
	t.Run("Methods()", testWithSpecFn(registry, Methods(http.MethodGet, http.MethodPost)))
	t.Run("Header()", testPredicatesWithoutSpec(Header("key", "value")))
	t.Run("HeaderRegexp()", testPredicatesWithoutSpec(HeaderRegexp("key", regex)))
	t.Run("Cookie()", testWithSpecFn(registry, Cookie("cookieName", regex)))
	t.Run("JWTPayloadAnyKV(single arg)", testWithSpecFn(registry, JWTPayloadAnyKV(kvPair)))
	t.Run("JWTPayloadAnyKV(multiple args)", testWithSpecFn(registry, JWTPayloadAnyKV(kvPair, kvPair)))
	t.Run("JWTPayloadAllKV(single arg)", testWithSpecFn(registry, JWTPayloadAllKV(kvPair)))
	t.Run("JWTPayloadAllKV(multiple args)", testWithSpecFn(registry, JWTPayloadAllKV(kvPair, kvPair)))
	t.Run("JWTPayloadAnyKVRegexp(single arg)", testWithSpecFn(registry, JWTPayloadAnyKVRegexp(kvRegexPair)))
	t.Run("JWTPayloadAnyKVRegexp(multiple args)", testWithSpecFn(registry, JWTPayloadAnyKVRegexp(kvRegexPair, kvRegexPair)))
	t.Run("JWTPayloadAllKVRegexp(single arg)", testWithSpecFn(registry, JWTPayloadAllKVRegexp(kvRegexPair)))
	t.Run("JWTPayloadAllKVRegexp(multiple args)", testWithSpecFn(registry, JWTPayloadAllKVRegexp(kvRegexPair, kvRegexPair)))
	t.Run("After()", testWithSpecFn(registry, After(time.Now())))
	t.Run("AfterWithDateString()", testWithSpecFn(registry, AfterWithDateString("2020-12-19T00:00:00+00:00")))
	t.Run("AfterWithUnixTime()", testWithSpecFn(registry, AfterWithUnixTime(time.Now().Unix())))
	t.Run("Before()", testWithSpecFn(registry, Before(time.Now())))
	t.Run("BeforeWithDateString()", testWithSpecFn(registry, BeforeWithDateString("2020-12-19T00:00:00+00:00")))
	t.Run("BeforeWithUnixTime()", testWithSpecFn(registry, BeforeWithUnixTime(time.Now().Unix())))
	t.Run("Between()", testWithSpecFn(registry, Between(time.Now(), time.Now().Add(time.Hour))))
	t.Run("BetweenWithDateString()", testWithSpecFn(registry, BetweenWithDateString("2020-12-19T00:00:00+00:00", "2020-12-19T01:00:00+00:00")))
	t.Run("BetweenWithUnixTime()", testWithSpecFn(registry, BetweenWithUnixTime(time.Now().Unix(), time.Now().Add(time.Hour).Unix())))
	t.Run("Cron()", testWithSpecFn(registry, Cron("* * * * *")))
	t.Run("QueryParam()", testWithSpecFn(registry, QueryParam("skipper")))
	t.Run("QueryParamWithValueRegex()", testWithSpecFn(registry, QueryParamWithValueRegex("skipper", regex)))
	t.Run("Source(single arg)", testWithSpecFn(registry, Source("127.0.0.1")))
	t.Run("Source(multiple args)", testWithSpecFn(registry, Source("127.0.0.1", "10.0.0.0/24")))
	t.Run("SourceFromLast(single arg)", testWithSpecFn(registry, SourceFromLast("127.0.0.1")))
	t.Run("SourceFromLast(multiple args)", testWithSpecFn(registry, SourceFromLast("127.0.0.1", "10.0.0.0/24")))
	t.Run("ClientIP(single arg)", testWithSpecFn(registry, ClientIP("127.0.0.1")))
	t.Run("ClientIP(multiple args)", testWithSpecFn(registry, ClientIP("127.0.0.1", "10.0.0.0/24")))
	t.Run("Tee()", testWithSpecFn(registry, Tee("skipper")))
	t.Run("Traffic()", testWithSpecFn(registry, Traffic(.25)))
	t.Run("TrafficSticky()", testWithSpecFn(registry, TrafficSticky(.25, "catalog-test", "A")))
}

func testWithSpecFn(registry Registry, predicate *eskip.Predicate) func(t *testing.T) {
	return func(t *testing.T) {
		if ps, ok := registry[predicate.Name]; ok {
			_, err := ps.Create(predicate.Args)
			if err != nil {
				t.Errorf("unexpected error while parsing %s predicate with args %s, %v", predicate.Name, predicate.Args, err)
			}
		} else {
			t.Errorf("predicate with name not found in registry, predicate=%s", predicate.Name)
		}
		//if predicateSpec.Name() != predicate.Name {
		//}
		//_, err := predicateSpec.Create(predicate.Args)
	}
}

func testPredicatesWithoutSpec(predicate *eskip.Predicate) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		switch predicate.Name {
		case HostRegexpName:
			_, err = ValidateHostRegexpPredicate(predicate)
		case PathRegexpName:
			_, err = ValidatePathRegexpPredicate(predicate)
		case MethodName:
			_, err = ValidateMethodPredicate(predicate)
		case HeaderName:
			_, err = ValidateHeaderPredicate(predicate)
		case HeaderRegexpName:
			_, err = ValidateHeaderRegexpPredicate(predicate)
		case WeightName:
			_, err = weight.ParseWeightPredicateArgs(predicate.Args)
		case PathName, PathSubtreeName:
			_, err = ProcessPathOrSubTree(predicate)
		default:
			t.Errorf("Unknown predicate provided %q", predicate.Name)
		}
		if err != nil {
			t.Errorf("unexpected error while parsing %s predicate with args %s, %v", predicate.Name, predicate.Args, err)
		}
	}
}
