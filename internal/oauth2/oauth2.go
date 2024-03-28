package oauth2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"time"

	mysql "goddd/pkg/oauth2-mysql"
	"goddd/pkg/token"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenInvalidScope = errors.New("invalid scope")
	ErrTokenUnregistered = errors.New("unregistered token")
)

type Server struct {
	logger            *logrus.Logger
	domain            string
	scope             string
	secretKey         string
	kid               string
	tokenExpiration   time.Duration
	server            *server.Server
	clientStore       *mysql.ClientStore
	tokenStore        *mysql.TokenStore
	jwtAccessGenerate *token.JWTAccessGenerate
}

func NewServer(logger *logrus.Logger, db *sqlx.DB) *Server {
	var (
		host          = "0.0.0.0"
		port          = "80"
		secretKey     = "secret"
		hours         = "10"
		expiryTime, _ = strconv.Atoi(hours)
		domain        = fmt.Sprintf("http://%s:%s", host, port)
		kid           = "alonsojl"
		scope         = "admin"

		tokenExpiration   = time.Duration(expiryTime) * time.Minute
		clientStore, _    = mysql.NewClientStore(db, mysql.WithClientStoreInitTableDisabled())
		tokenStore, _     = mysql.NewTokenStore(db, mysql.WithTokenStoreInitTableDisabled())
		jwtAccessGenerate = token.NewJWTAccessGenerate(kid, []byte(secretKey), jwt.SigningMethodHS512)
	)
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(clientStore)
	manager.MapTokenStorage(tokenStore)
	manager.MapAccessGenerate(jwtAccessGenerate)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(false)

	return &Server{
		logger:            logger,
		domain:            domain,
		scope:             scope,
		secretKey:         secretKey,
		kid:               kid,
		tokenExpiration:   tokenExpiration,
		server:            srv,
		clientStore:       clientStore,
		tokenStore:        tokenStore,
		jwtAccessGenerate: jwtAccessGenerate,
	}
}

func (a *Server) GetCredentials(clientId string) (oauth2.ClientInfo, error) {
	clientCredentials, err := a.clientStore.GetByID(context.Background(), clientId)
	if err != nil {
		a.logger.WithError(err).Error("error getting client credentials")
		return nil, err
	}
	return clientCredentials, nil
}

func (a *Server) CreateCredentials() (oauth2.ClientInfo, error) {
	var (
		clientID     = uuid.New().String()
		clientSecret = uuid.New().String()
		userID       = uuid.New().String()
	)

	clientCredentials := &models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: a.domain,
		UserID: userID,
	}

	err := a.clientStore.Create(context.Background(), clientCredentials)
	if err != nil {
		a.logger.WithError(err).Error("error saving client credentials")
		return &models.Client{}, err
	}

	return clientCredentials, nil
}

func (a *Server) CreateAccessToken(clientCredentials oauth2.ClientInfo) (string, error) {
	createAt := time.Now()
	tokenInfo := &models.Token{
		ClientID:        clientCredentials.GetID(),
		UserID:          clientCredentials.GetUserID(),
		Scope:           a.scope,
		AccessCreateAt:  createAt,
		AccessExpiresIn: a.tokenExpiration,
	}

	accessToken, _, err := a.jwtAccessGenerate.Token(
		context.Background(),
		&oauth2.GenerateBasic{
			Client:    clientCredentials,
			CreateAt:  createAt,
			TokenInfo: tokenInfo,
		},
		false,
	)

	if err != nil {
		a.logger.WithError(err).Error("error generating access token")
		return "", err
	}

	tokenInfo.Access = accessToken
	err = a.tokenStore.Create(context.Background(), tokenInfo)
	if err != nil {
		a.logger.WithError(err).Error("error saving access token")
		return "", err
	}

	return accessToken, nil
}

func (a *Server) ValidateToken(r *http.Request) error {
	var (
		ctx         = r.Context()
		accessToken = r.Header.Get("Authorization")
	)

	tokenStr, err := a.jwtAccessGenerate.Validate(accessToken)
	if err != nil {
		a.logger.WithError(err).Error("error validating token access")
		return err
	}

	tokenInfo, err := a.tokenStore.GetByAccess(ctx, tokenStr)
	if err != nil {
		a.logger.WithError(err).Error("error getting token info")
		return err
	}

	if tokenInfo == nil {
		err = ErrTokenUnregistered
		a.logger.WithError(err).Error("error validating token info")
		return err
	}

	var (
		scopes       = strings.Split(tokenInfo.GetScope(), ",")
		isValidScope = false
	)

	for _, scope := range scopes {
		if scope == a.scope {
			isValidScope = true
		}
	}

	if !isValidScope {
		err = ErrTokenInvalidScope
		a.logger.WithError(err).Error("error validating token scope")
		return err
	}

	return nil
}

// func (a *Server) RevokeTokens(ctx context.Context, clientID string) error {
// 	return a.tokenStore.RemoveByClientID(ctx, clientID)
// }

func (a *Server) ValidationBearerToken(r *http.Request) error {
	tokenInfo, err := a.server.ValidationBearerToken(r)
	if err != nil {
		a.logger.WithError(err).Error("error validating token access")
		return err
	}

	var (
		scopes       = strings.Split(tokenInfo.GetScope(), ",")
		isValidScope = false
	)

	for _, scope := range scopes {
		if scope == a.scope {
			isValidScope = true
		}
	}

	if !isValidScope {
		err = ErrTokenInvalidScope
		a.logger.WithError(err).Error("error validating token scope")
		return err
	}

	return nil
}
