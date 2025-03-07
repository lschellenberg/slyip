package verifier

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	RoleAdmin = "admin"
	RoleBasic = "basic"
)

// Principal is a simple struct containing the id and scopes
// Can be used to verify simply the account identity without having to read the full token
type Principal struct {
	ID               string
	ECDSAAddress     string
	SLYWalletAddress string
	Role             string
	Scopes           []string
	Audiences        []string
}

func (p Principal) IsAdmin() bool {
	return p.Role == RoleAdmin
}

// ParseAuthorizationBearer expects a header Authorization to contain Bearer <token>
// <token> is returned if parsing is correct
// a BadRequest is returned if parsing is not possible
func ParseAuthorizationBearer(req *http.Request) (string, error) {
	val := req.Header.Get("Authorization")

	if val == "" {
		return "", fmt.Errorf("check authorization header")
	}

	parts := strings.SplitN(val, " ", 2)

	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf("malformed bearer authorization")
	}

	return parts[1], nil
}

func (p Principal) HasPermission(permission string) bool {
	for _, s := range p.Scopes {
		if permission == s {
			return true
		}
	}
	return false
}
