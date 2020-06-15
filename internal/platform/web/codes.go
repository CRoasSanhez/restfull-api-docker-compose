package web

import "github.com/CRoasSanhez/yofio-test/internal/responses"

// User Response Codes are used for a deep description of the error
var (
	EmailRepeated             = responses.NewErrorCode(1000, "Email already in use")
	PhoneRepeated             = responses.NewErrorCode(1000, "Phone already in use")
	EmailInvalid              = responses.NewErrorCode(1001, "Email invalid")
	PhoneInvalid              = responses.NewErrorCode(1001, "Phone invalid")
	NameInvalid               = responses.NewErrorCode(1001, "Name invalid")
	PasswordInvalid           = responses.NewErrorCode(1001, "Password invalid")
	EmailNotFound             = responses.NewErrorCode(1002, "Email not found")
	PhoneNotFound             = responses.NewErrorCode(1002, "Phone not found")
	ErrorInvalidData          = responses.NewErrorCode(1003, "Error invalid data")
	UserBlocked               = responses.NewErrorCode(1005, "User blocked")
	MaxLoginAttemptsReached   = responses.NewErrorCode(1006, "Max number of login attempts reached")
	MaxPaymentAttemptsReached = responses.NewErrorCode(1007, "Max number of payment attempts reached")
	UnsuccessPayment          = responses.NewErrorCode(1008, "Unsuccess payment")
)

// Server Error Response Codes are used for a deep description of the error
var (
	ErrorDBConnection       = responses.NewErrorCode(3000, "DB connection error")
	ErrorGenerateJWT        = responses.NewErrorCode(3005, "generating JWT")
	ErrorSigningJWT         = responses.NewErrorCode(3006, "Signing JWT")
	ErrorValidatingJWT      = responses.NewErrorCode(3007, "Validating JWT")
	ErrorRegisterUser       = responses.NewErrorCode(3003, "User registratin")
	ErrorLoginUser          = responses.NewErrorCode(3003, "User not login")
	ErrorSavingAttempts     = responses.NewErrorCode(3008, "Saving Login Attemprs")
	ErrorBlockingUser       = responses.NewErrorCode(3009, "Blocking user")
	ErrorUpdatingMembership = responses.NewErrorCode(3010, "Updating membership")
	ErrorPaymentsConsult    = responses.NewErrorCode(3011, "Retreiving membership payments consult")
)
