package pageHandlers

import (
	"context"
	"log/slog"
	"net/mail"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/alekLukanen/go-templ-htmx-example-app/handlers"
	"github.com/alekLukanen/go-templ-htmx-example-app/services"
	"github.com/alekLukanen/go-templ-htmx-example-app/ui/pages"
)

type SignupPageHandler struct {
	ctx    context.Context
	logger *slog.Logger

	*services.ServiceMesh
}

func NewSignupPageHandler(ctx context.Context, logger *slog.Logger, serviceMesh *services.ServiceMesh) *SignupPageHandler {
	return &SignupPageHandler{ctx: ctx, logger: logger, ServiceMesh: serviceMesh}
}

func (obj *SignupPageHandler) RegisterPublicRoutes(echoHandler *echo.Echo) {
	echoHandler.GET("/signup", obj.BasePage)
	echoHandler.POST("/signup/validate-email", obj.ValidateEmail)
	echoHandler.POST("/signup/validate-passwords", obj.ValidatePasswords)
	echoHandler.POST("/signup/submit", obj.Submit)
}

func (obj *SignupPageHandler) BasePage(echoCtx echo.Context) error {
	page := pages.SignupPage(pages.SignupInputFormParams{})
	handlers.Render(echoCtx, &page)
	return nil
}

func (obj *SignupPageHandler) FormToParams(echoCtx echo.Context) pages.SignupInputFormParams {
	email := echoCtx.FormValue("email")
	password1 := echoCtx.FormValue("password1")
	password2 := echoCtx.FormValue("password2")

	showInvalidPasswordFlag := false
	showNonMatchingPasswordFlag := false
	showInvalidEmailFlag := false
	showTakenEmailFlag := false
	submitButtonDisabled := false

	_, emailParseErr := mail.ParseAddress(email)
	if emailParseErr != nil {
		showInvalidEmailFlag = true
		submitButtonDisabled = true
	}

	if !obj.UserAuthenticationService.PasswordIsValid(password1) {
		showInvalidPasswordFlag = true
	}

	if len(password1) > 0 && password1 != password2 {
		showNonMatchingPasswordFlag = true
	}

	if (len(password1) == 0 || len(password2) == 0) || showInvalidPasswordFlag || showNonMatchingPasswordFlag {
		submitButtonDisabled = true
	}

	params := pages.SignupInputFormParams{
		Email:                       email,
		Password1:                   password1,
		Password2:                   password2,
		ShowInvalidPasswordFlag:     showInvalidPasswordFlag,
		ShowNonMatchingPasswordFlag: showNonMatchingPasswordFlag,
		ShowInvalidEmailFlag:        showInvalidEmailFlag,
		ShowTakenEmailFlag:          showTakenEmailFlag,
		SubmitButtonDisabled:        submitButtonDisabled,
	}

	return params
}

func (obj *SignupPageHandler) ValidateEmail(echoCtx echo.Context) error {
	params := obj.FormToParams(echoCtx)
	emailErrors := pages.EmailErrors(params)
	button := pages.SignupPageButton(params, "outerHTML")

	handlers.Render(echoCtx, &emailErrors)
	handlers.Render(echoCtx, &button)
	return nil
}

func (obj *SignupPageHandler) ValidatePasswords(echoCtx echo.Context) error {
	params := obj.FormToParams(echoCtx)
	passwordErrors := pages.PasswordErrors(params)
	button := pages.SignupPageButton(params, "outerHTML")

	handlers.Render(echoCtx, &passwordErrors)
	handlers.Render(echoCtx, &button)
	return nil
}

func (obj *SignupPageHandler) Submit(echoCtx echo.Context) error {

	params := obj.FormToParams(echoCtx)

	signupCompleted := false
	if params.FormAppearsValid() {

		email := echoCtx.FormValue("email")
		password := echoCtx.FormValue("password1")
		_, err := obj.UserAuthenticationService.Signup(echoCtx.Request().Context(), email, password)

		if err == services.ErrEmailTaken {
			params.ShowTakenEmailFlag = true
			params.SubmitButtonDisabled = true

		} else if err == services.ErrInvalidPassword {
			params.ShowInvalidPasswordFlag = true
			params.SubmitButtonDisabled = true

		} else if err != nil {
			params.ShowInvalidPasswordFlag = true
			params.SubmitButtonDisabled = true

		} else {
			signupCompleted = true
		}

	}

	var component templ.Component
	if signupCompleted {
		component = pages.SignupSuccess()

	} else {
		component = pages.SignupInputForm(params)
	}
	handlers.Render(echoCtx, &component)
	return nil
}
