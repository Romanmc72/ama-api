package endpoints_test

import (
	"ama/api/constants"
	"ama/api/endpoints"
	"ama/api/test"
	"ama/api/test/fixtures"
	"slices"
	"testing"
)

func TestGetQueryParamToIntAndInt64(t *testing.T) {
	paramName := "how-many-bacon-and-eggs-would-you-like?"
	var defaultValuei64 int64 = 10
	defaultValue := int(defaultValuei64)
	testCases := []struct {
		name string
		ctx  *test.MockAPIContext
		want int64
	}{
		{
			name: "Success - Default",
			ctx: &test.MockAPIContext{
				QueryValues: map[string][]string{},
			},
			want: defaultValuei64,
		},
		{
			name: "Success - Default On Invalid Input",
			ctx: &test.MockAPIContext{
				QueryValues: map[string][]string{
					paramName: {"All of the bacon and eggs you have."},
				},
			},
			want: defaultValuei64,
		},
		{
			name: "Success - Retrieved First",
			ctx: &test.MockAPIContext{
				QueryValues: map[string][]string{
					paramName: {"21", "100"},
				},
			},
			want: 21,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := endpoints.GetQueryParamToInt(tc.ctx, paramName, defaultValue)
			if int(tc.want) != got {
				t.Errorf("GetQueryParamToInt() wanted = %d; got = %d", tc.want, got)
			}
		})
		t.Run(tc.name+"(int64)", func(t *testing.T) {
			got := endpoints.GetQueryParamToInt64(tc.ctx, paramName, defaultValuei64)
			if tc.want != got {
				t.Errorf("GetQueryParamToInt64() wanted = %d; got = %d", tc.want, got)
			}
		})
	}
}

func TestGetReadQuestionsParamsWithDefaults(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         *test.MockAPIContext
		isRandom    bool
		wantLimit   int
		wantFinalId string
		wantTags    []string
	}{
		{
			name: "Success - Defaults",
			ctx: &test.MockAPIContext{
				QueryValues: map[string][]string{},
			},
			isRandom:    false,
			wantLimit:   constants.DefaultLimit,
			wantFinalId: "",
			wantTags:    []string{},
		},
		{
			name: "Success - Random finalId",
			ctx: &test.MockAPIContext{
				QueryValues: map[string][]string{
					constants.LimitParam:  {"10"},
					constants.RandomParam: {"true"},
					constants.TagParam:    {"a", "b", "c"},
				},
			},
			isRandom:    true,
			wantLimit:   10,
			wantFinalId: "<ignore, should be random>",
			wantTags:    []string{"a", "b", "c"},
		},
		{
			name: "Success - Pre Set finalId",
			ctx: &test.MockAPIContext{
				QueryValues: map[string][]string{
					constants.FinalIdParam: {fixtures.NewId},
				},
			},
			isRandom:    false,
			wantLimit:   constants.DefaultLimit,
			wantFinalId: fixtures.NewId,
			wantTags:    []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLimit, gotFinalId, gotTags := endpoints.GetReadQuestionsParamsWithDefaults(tc.ctx)

			if tc.wantLimit != gotLimit || (!tc.isRandom && tc.wantFinalId != gotFinalId) || !slices.Equal(tc.wantTags, gotTags) {
				t.Errorf(
					"GetReadQuestionsParamsWithDefaults() "+
						"wantLimit = %d; gotLimit = %d; "+
						"wantFinalId = %s (randomFinalId = %v); gotFinalId = %s; "+
						"wantTags = %v; gotTags = %v",
					tc.wantLimit,
					gotLimit,
					tc.wantFinalId,
					tc.isRandom,
					gotFinalId,
					tc.wantTags,
					gotTags,
				)
			}
		})
	}
}
