package jwt

import (
	"net/http"
	"reflect"
	"testing"
)

func init() {
	ValidateOnParse(false)
}

func TestClaims(t *testing.T) {

	const usertoken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE0NzY0ODAwLCJleHAiOjE1MTc0NDMyMDAsImVtYWlsIjoiSm9objEyM0BnbWFpbC5jb20iLCJsYXN0bmFtZSI6IkRvZSIsImZpcnN0bmFtZSI6IkpvaG4iLCJhdWQiOiJleGFtcGxlLmNvbSJ9.JJ_V8UPs2WX6OXEvw44zj0hcnaRr425A66boJZOybzE"
	req, _ := http.NewRequest("GET", "/echo", nil)
	req.Header.Set("Authorization", usertoken)
	claims, _ := Claims(req)

	tests := []struct {
		claim    string
		expected interface{}
	}{
		{"sub", "1234567890"},
		{"iat", float64(1514764800)},
		{"exp", float64(1517443200)},
		{"email", "John123@gmail.com"},
		{"lastname", "Doe"},
		{"firstname", "John"},
		{"aud", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.claim, func(t *testing.T) {
			actual := claims.Get(tt.claim).value
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("GetClaim() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestClaimsInvalidToken(t *testing.T) {
	tests := []struct {
		token string
	}{
		{""},
		{"Bearer "},
		{"Bearer foo"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/echo", nil)
			req.Header.Set("Authorization", tt.token)
			_, err := Claims(req)
			if err == nil {
				t.Error("Expected error, got none")
			}
		})
	}
}

func TestClaimsType(t *testing.T) {
	const (
		asString = iota
		asFloat64
		asInt64
	)
	tests := []struct {
		name      string
		input     interface{}
		converter int
		expected  interface{}
	}{
		{"string to string", "foo", asString, "foo"},
		{"float64 to string", float64(42), asString, "42"},
		{"float64 to float64", float64(42), asFloat64, float64(42)},
		{"float64 to int64", float64(42), asInt64, int64(42)},
		{"int64 to int64", int64(42), asInt64, int64(42)},
		{"int64 to string", int64(42), asString, "42"},
		{"nil to string", nil, asString, ""},
		{"nil to float64", nil, asFloat64, float64(0)},
		{"nil to int64", nil, asInt64, int64(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual interface{}
			ct := claimType{tt.input}
			switch tt.converter {
			case asString:
				actual = ct.AsString()
			case asFloat64:
				actual = ct.AsFloat64()
			case asInt64:
				actual = ct.AsInt64()
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("claimType = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func Test_claims_Get(t *testing.T) {
	tests := []struct {
		name  string
		i     *claims
		claim string
		want  *claimType
	}{
		{"nil claims", nil, "sub", &claimType{}},
		{"empty claims", &claims{}, "sub", &claimType{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Get(tt.claim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("claims.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
