package jwt

import (
	"encoding/base64"
	"errors"
	"github.com/fangbc5/gogo/core/auth"
	"github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

var (
	// ErrNotFound is returned when a token cannot be found.
	ErrNotFound = errors.New("token not found")
	// ErrEncodingToken is returned when the service encounters an error during encoding.
	ErrEncodingToken = errors.New("error encoding the token")
	// ErrInvalidToken is returned when the token provided is not valid.
	ErrInvalidToken = errors.New("invalid token provided")
)

// authClaims to be encoded in the JWT.
type authClaims struct {
	Type     string            `json:"type"`
	Scopes   []string          `json:"scopes"`
	Metadata map[string]string `json:"metadata"`

	jwt.RegisteredClaims
}

type jwtImpl struct {
	sync.Mutex
	options auth.Options
}

// NewAuth returns a new instance of the Auth service.
func NewAuth(opts ...auth.Option) auth.Auth {
	j := new(jwtImpl)
	j.Init(opts...)
	return j
}

func NewRules() auth.Rules {
	return new(jwtRules)
}

func (j *jwtImpl) Init(opts ...auth.Option) {
	j.Lock()
	defer j.Unlock()

	for _, o := range opts {
		o(&j.options)
	}
}

func (j *jwtImpl) Options() auth.Options {
	j.Lock()
	defer j.Unlock()
	return j.options
}

func (j *jwtImpl) Generate(id string, opts ...auth.GenerateOption) (*auth.Token, error) {
	options := auth.NewGenerateOptions(opts...)
	account := &auth.Account{
		ID:       id,
		Type:     options.Type,
		Scopes:   options.Scopes,
		Metadata: options.Metadata,
		Issuer:   j.Options().Namespace,
	}

	// generate a JWT secret which can be provided to the Token() method
	// and exchanged for an access token
	access, err := j.generate(account, auth.WithExpiry(options.Expiry))
	if err != nil {
		return nil, err
	}

	refresh, err := j.generate(account, auth.WithExpiry(options.Expiry+time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	return &auth.Token{
		Created:      access.Created,
		Expiry:       access.Expiry,
		AccessToken:  access.AccessToken,
		RefreshToken: refresh.AccessToken,
	}, nil
}

// generate a new JWT.
func (j *jwtImpl) generate(acc *auth.Account, opts ...auth.GenerateOption) (*auth.Token, error) {
	// decode the private key
	priv, err := base64.StdEncoding.DecodeString(j.options.PrivateKey)
	if err != nil {
		return nil, err
	}

	// parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, ErrEncodingToken
	}

	// parse the options
	options := auth.NewGenerateOptions(opts...)

	// generate the JWT
	expiry := time.Now().Add(options.Expiry)
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, authClaims{
		acc.Type, acc.Scopes, acc.Metadata, jwt.RegisteredClaims{
			Subject:   acc.ID,
			Issuer:    acc.Issuer,
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	})
	tok, err := t.SignedString(key)
	if err != nil {
		return nil, err
	}

	// return the token
	return &auth.Token{
		AccessToken: tok,
		Expiry:      expiry,
		Created:     time.Now(),
	}, nil
}

func (j *jwtImpl) Inspect(t string) (*auth.Account, error) {
	// decode the public key
	pub, err := base64.StdEncoding.DecodeString(j.options.PublicKey)
	if err != nil {
		return nil, err
	}

	// parse the public key
	res, err := jwt.ParseWithClaims(t, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(pub)
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	// validate the token
	if !res.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := res.Claims.(*authClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// return the token
	return &auth.Account{
		ID:       claims.Subject,
		Issuer:   claims.Issuer,
		Type:     claims.Type,
		Scopes:   claims.Scopes,
		Metadata: claims.Metadata,
	}, nil
}

func (j *jwtImpl) Refresh(opts ...auth.GenerateOption) (*auth.Token, error) {
	options := auth.NewGenerateOptions(opts...)

	rt := options.RefreshToken

	account, err := j.Inspect(rt)
	if err != nil {
		return nil, err
	}

	access, err := j.generate(account, auth.WithExpiry(options.Expiry))
	if err != nil {
		return nil, err
	}

	return &auth.Token{
		Created:      access.Created,
		Expiry:       access.Expiry,
		AccessToken:  access.AccessToken,
		RefreshToken: rt,
	}, nil
}

func (j *jwtImpl) String() string {
	return "jwt"
}

type jwtRules struct {
	sync.Mutex
	rules []*auth.Rule
}

func (j *jwtRules) Grant(rule *auth.Rule) error {
	j.Lock()
	defer j.Unlock()
	j.rules = append(j.rules, rule)
	return nil
}

func (j *jwtRules) Revoke(rule *auth.Rule) error {
	j.Lock()
	defer j.Unlock()

	rules := make([]*auth.Rule, 0, len(j.rules))
	for _, r := range j.rules {
		if r.ID != rule.ID {
			rules = append(rules, r)
		}
	}

	j.rules = rules
	return nil
}

func (j *jwtRules) Verify(acc *auth.Account, res *auth.Resource, opts ...auth.VerifyOption) error {
	j.Lock()
	defer j.Unlock()

	var options auth.VerifyOptions
	for _, o := range opts {
		o(&options)
	}

	return auth.Verify(j.rules, acc, res)
}

func (j *jwtRules) List(opts ...auth.ListOption) ([]*auth.Rule, error) {
	j.Lock()
	defer j.Unlock()
	return j.rules, nil
}
