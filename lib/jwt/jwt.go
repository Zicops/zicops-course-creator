package jwt

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	jwt "github.com/golang-jwt/jwt"
)

// ValidateOnParse will say whether the token should be validated when the token is parsed
// Default is yes
func ValidateOnParse(yes bool) {
	validate = yes
}

var validate = false

type claims struct {
	allclaims map[string]interface{}
}

// Claims returns the claims in the request
func Claims(r *http.Request) (*claims, error) {
	return ClaimsFromToken(GetToken(r))
}

// ClaimsFromToken creates claims from given token
func ClaimsFromToken(token string) (*claims, error) {
	allclaims, e := GetClaims(token)
	if e != nil {
		return nil, e
	}
	return &claims{allclaims}, nil
}

// Get returns a claim
func (i *claims) Get(claim string) *claimType {
	ct := &claimType{}
	if i == nil {
		return ct
	}
	if value, ok := i.allclaims[claim]; ok {
		ct.value = value
	}
	return ct
}

// GetToken makes a token from a request
func GetToken(r *http.Request) string {
	authorization := r.Header.Get("authorization")
	if authorization != "" {
		bearer := "Bearer "
		if string([]rune(authorization)[:len(bearer)]) == bearer {
			return string([]rune(authorization)[len(bearer):])
		}

		return authorization
	}
	return ""
}

// GetTokenWebsocket makes a token from a request
func GetTokenWebsocket(r *http.Request) string {
	authorization := r.Header.Get("sec-websocket-protocol")
	if authorization == "" {
		authorization = r.Header.Get("Sec-WebSocket-Protocol")
	}
	if authorization != "" {
		return authorization
	}
	return ""
}

func GetClaims(stoken string) (map[string]interface{}, error) {
	token, _ := jwt.Parse(stoken, func(t *jwt.Token) (interface{}, error) {
		return t, nil
	})

	if token == nil {
		return nil, fmt.Errorf("Unable to get claims from token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unable to get claims from token")
	}
	if validate && !token.Valid {
		return nil, fmt.Errorf("not valid token")
	}

	return claims, nil
}

type claimType struct {
	value interface{}
}

// AsString casts to string
func (ct *claimType) AsString() string {
	if ct.value != nil {
		switch ct.value.(type) {
		case string:
			return ct.value.(string)
		case float64:
			return strconv.FormatFloat(ct.value.(float64), 'g', 6, 64)
		case int64:
			return strconv.FormatInt(ct.value.(int64), 10)
		}
	}
	return ""
}

// AsFloat64 casts to float64
func (ct *claimType) AsFloat64() float64 {
	if ct.value != nil {
		return ct.value.(float64)
	}
	return 0
}

// AsInt64 converts numbers to int64
func (ct *claimType) AsInt64() int64 {
	if ct.value != nil {
		switch reflect.TypeOf(ct.value).Kind() {
		case reflect.Float64:
			return int64(ct.value.(float64))
		case reflect.Int64:
			return ct.value.(int64)
		}
	}
	return 0
}

// GetUserID get a user id from request, will be subid if it exists, otherwise sub
func GetUserID(r *http.Request) (string, error) {

	claims, err := Claims(r)

	if err != nil {
		return "", err
	}

	subid := claims.Get("subid").AsString()
	if subid != "" {
		return subid, nil
	}

	return claims.Get("sub").AsString(), nil
}

// GetUserEmail get user email address from request while registering action/user it is useful
func GetUserEmail(r *http.Request) (string, error) {
	claims, err := Claims(r)

	if err != nil {
		return "", err
	}
	userid := claims.Get("userid").AsString()
	if userid != "" {
		return userid, nil
	}
	return claims.Get("email").AsString(), nil
}

// GetAudience get user email address from request while registering action/user it is useful
func GetAudience(r *http.Request) (string, error) {
	claims, err := Claims(r)

	if err != nil {
		return "", err
	}
	aud := claims.Get("aud").AsString()
	if aud != "" {
		return aud, nil
	}
	return "everyone", nil
}
