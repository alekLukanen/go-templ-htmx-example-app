package services

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/alekLukanen/go-templ-htmx-example-app/core/database/queries"
	"github.com/alekLukanen/go-templ-htmx-example-app/core/settings"
)

var ErrInvalidPassword error = errors.New("ErrInvalidPassword")
var ErrEmailTaken error = errors.New("ErrEmailTaken")
var ErrPasswordIncorrect error = errors.New("ErrPasswordIncorrect")
var ErrParsingJWTToken error = errors.New("ErrParsingJWTToken")
var ErrInvalidToken error = errors.New("ErrInvalidToken")

const TOKEN_DURATION time.Duration = time.Minute * 30

type JWTScopeClaims struct {
	UserID int32  `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type UserAuthenticationService struct {
	ctx       context.Context
	logger    *slog.Logger
	dbQueries *queries.Queries
}

func NewUserAuthenticationService(ctx context.Context, logger *slog.Logger, dbQueries *queries.Queries) UserAuthenticationService {
	return UserAuthenticationService{
		ctx:       ctx,
		logger:    logger,
		dbQueries: dbQueries,
	}
}

func (obj *UserAuthenticationService) Signup(ctx context.Context, email string, password string) (*queries.User, error) {
	if !obj.EmailIsValid(email) || !obj.PasswordIsValid(password) {
		return nil, ErrInvalidPassword
	}

	if taken, err := obj.dbQueries.EmailTaken(ctx, email); err != nil {
		return nil, err
	} else if taken {
		return nil, ErrEmailTaken
	}

	hashedPassword, err := obj.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := obj.dbQueries.CreateUser(ctx, queries.CreateUserParams{Email: email, Password: hashedPassword, Enabled: true})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (obj *UserAuthenticationService) Signin(ctx context.Context, email, password string) (string, error) {
	// implementation
	user, err := obj.dbQueries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !obj.PasswordMatchesUserPassword(user.Password, password) {
		return "", ErrPasswordIncorrect
	}

	token, err := obj.GenerateJWT(&user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (obj *UserAuthenticationService) UserFromEchoContext(echoCtx echo.Context) (*queries.User, error) {

	tokenUser := echoCtx.Get("user").(*jwt.Token)
	claims := tokenUser.Claims.(*JWTScopeClaims)
	userId := claims.UserID
	user, err := obj.dbQueries.GetUserById(echoCtx.Request().Context(), userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (obj *UserAuthenticationService) GenerateJWT(user *queries.User) (string, error) {
	var mySigningKey = []byte(settings.JWT_SECRET_KEY)

	claims := &JWTScopeClaims{
		user.ID,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(TOKEN_DURATION)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenValue, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenValue, nil
}

func (obj *UserAuthenticationService) EmailIsValid(email string) bool {
	_, emailParseErr := mail.ParseAddress(email)
	if emailParseErr != nil {
		return false
	}
	return true
}

func (obj *UserAuthenticationService) PasswordIsValid(password string) bool {
	if len(password) >= 8 && len(password) <= 25 {
		return true
	}
	return false
}

func (obj *UserAuthenticationService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (obj *UserAuthenticationService) PasswordMatchesUserPassword(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
