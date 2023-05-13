package common

const (
	// ERROR CODES and MESSAGES
	WRONG_LOGIN_REQUEST_CODE     = 4000
	REFRESH_TOKEN_INCORRECT_CODE = 4001
	LOGIN_FAILED_CODE            = 4002
	LOGOUT_FAILED_CODE           = 4003
	REGISTER_FAILED_CODE         = 4004
	UPDATE_PASSWORD_FAILED_CODE  = 4005
	VERIFICATION_FAILED_CODE     = 4006
	SUCCESS_CODE                 = 2000
	REFRESH_TOKEN_INCORRECT      = "Refresh token is incorrect"
	VERIFICATION_FAILED          = "Token is Invalid"
	WRONG_PASSWORD               = "incorrect password"
	INVALID_JSON                 = "Invalid json provided"
	LOGOUT_FAILED                = "Logout failed"
	TOKEN_CREATION_FAILED        = "Failed to create token"
	WRONG_PAYLOAD_TYPE           = "Payload type should be email, wechat or phone"
	USER_EXISTS                  = "there is an account with this email, try a different email"
	REGISTRATION_FAILED          = "sorry, the email is already in use, try a different email"
	NO_RECORD_FOUND              = "sorry, email does not exist"
	PROFILE_QUERY_FAILED         = "this profile does not exist, please register"
	PROFILE_CREATION_FAILED      = "failed to create profile, please try again"
	SESSION_CREATION_FAILED      = "failed to create a session, try again"
	PUSH_REDIS_FAILED            = "Failed to push token into redis"
	LOGIN_SUCCESS                = "Login successfully"
	LOGOUT_SUCCESS               = "Logout successfully"
	REGISTER_SUCCESS             = "User registered successfully"
	REFRESH_SUCCESS              = "Token refreshed successfully"
	PASSWORD_CHANGED             = "Password Changed successfully"
	VERIFICATION_OK              = "Token refreshed successfully"
	UNSUCCESSFUL                 = "Unsuccessfully"
	PASSWORD_RESET_FAILED        = "failed to reset password"
	EMAIL_NOT_MATCHED            = "email does not match. enter correct email"
)
