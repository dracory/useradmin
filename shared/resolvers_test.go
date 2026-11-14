package shared_test

import (
	"testing"

	"github.com/dracory/useradmin/shared"
)

func TestSearchOpConstants(t *testing.T) {
	if shared.SearchOpEquals != "eq" ||
		shared.SearchOpContains != "contains" ||
		shared.SearchOpNotContains != "not_contains" ||
		shared.SearchOpStartsWith != "starts_with" {
		t.Fatal("unexpected search ops")
	}
	if shared.SearchFieldEmail != "email" ||
		shared.SearchFieldFirstName != "first_name" ||
		shared.SearchFieldLastName != "last_name" ||
		shared.SearchFieldPhone != "phone" ||
		shared.SearchFieldBusinessName != "business_name" {
		t.Fatal("unexpected search fields")
	}
	if shared.SearchAnd != "AND" || shared.SearchOr != "OR" {
		t.Fatal("unexpected combine ops")
	}
}

func TestSearchConditionStruct(t *testing.T) {
	c := shared.SearchCondition{
		Field:       shared.SearchFieldEmail,
		Op:          shared.SearchOpContains,
		Value:       "@t.com",
		CombineWith: shared.SearchAnd,
	}
	if c.Field != shared.SearchFieldEmail || c.CombineWith != shared.SearchAnd {
		t.Fatal("unexpected fields")
	}
}

func TestCountryTimezoneStructs(t *testing.T) {
	c := shared.Country{IsoCode2: "US", Name: "United States"}
	tz := shared.Timezone{Code: "UTC"}
	if c.IsoCode2 != "US" || c.Name == "" || tz.Code != "UTC" {
		t.Fatal("unexpected fields")
	}
}
