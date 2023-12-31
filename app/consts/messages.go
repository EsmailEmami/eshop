package consts

const (
	// Messages
	Created              = "New record created successfully"
	Updated              = "Update performed successfully"
	Deleted              = "Record deleted successfully"
	VerifivationCodeSent = "Verification code sent successfully."
	RegistrationDone     = "Registration completed successfully."
	LoggedOut            = "Logged out successfully."
	OperationDone        = "Operation completed successfully"

	// Errors
	CacheRecordNotFound              = "Cache Record Not Found"
	BadRequest                       = "Invalid data provided"
	InvalidRequest                   = "Invalid request"
	InternalServerError              = "An internal server error occurred"
	UnauthorizedError                = "Authentication error"
	ForbiddenError                   = "Unauthorized access to the requested resource."
	ValidationError                  = "Data validation error"
	Required                         = "This field cannot be empty."
	InvalidValue                     = "Invalid value entered."
	RecordNotFound                   = "Requested record not found."
	FileNotFound                     = "Requested file not found."
	InvalidUsername                  = "Invalid username entered."
	InvalidUsernameCharacters        = "Invalid username entered. Only lowercase Latin letters and numbers are allowed."
	UsernameAlreadyExists            = "The entered username already exists."
	MobileAlreadyExists              = "The entered mobile number is already in use."
	UnableToCreate                   = "An error occurred while creating the record."
	UnableToUpdate                   = "An error occurred while updating the record."
	LoginFailed                      = "Invalid username or password entered."
	TokenError_MissingUserID         = "Invalid token. User ID is missing."
	TokenError_InvalidUserID         = "Invalid token. User ID is invalid."
	InvalidVoiceFileMimeType         = "Invalid audio file format selected."
	InvalidPdfFileMimeType           = "Invalid PDF file format selected."
	InvalidExcelFileMimeType         = "Invalid Excel file format selected."
	InvalidSvgFileMimeType           = "Invalid SVG file format selected."
	InvalidFileMimeType              = "Invalid file format selected."
	InvalidEndDate                   = "Invalid end date entered."
	InvalidStartDate                 = "Invalid start date entered."
	UnableToSendVerificationCode     = "An error occurred while sending the verification code. Please try again later."
	InvalidMobileNumber              = "Invalid mobile number entered."
	InvalidCharacters                = "The entered value contains invalid characters."
	InvalidColor                     = "The entered value is not a valid color code."
	HasActiveVerificationCodeRequest = "An activation code has already been sent to you. Please try again later."
	PasswordIsShort                  = "The password must be at least 8 characters long and include lowercase letters, uppercase letters, and special characters."
	InvalidVerificationCode          = "The entered verification code is invalid."
	InvalidPeyvastFile               = "An error occurred while validating the attachment file."
	InvalidPeyvastContentType        = "The uploaded file type is not allowed."
	InvalidFileHeader                = "The submitted file is not valid."
	InvalidFileContentType           = "The uploaded file type is not allowed."
	InvalidNumeric                   = "The entered value must be numeric."
	UnableToChangeSystemData         = "Cannot make changes to system information."
	InvalidTime                      = "Invalid time entered."
	InvalidPhone                     = "Invalid phone number entered."
	UserIsDisabled                   = "The user account is disabled."
	UnableToUpdateField              = "Cannot update the selected field."
	InvalidDate                      = "Invalid date entered."
	MinIsZero                        = "The minimum acceptable value is 0."
	MaxIsOne                         = "The maximum acceptable value is 1."
	InvalidGeoLocation               = "The submitted geographical location is invalid."
	ExistedCode                      = "The entered code has already been registered."
	ExistedTitle                     = "The entered title has already been registered."
	InvalidFileSize                  = "The uploaded file size is not allowed."
	OrderAlreadyRegistered           = "The entered order has already been registered."
	InvalidPriority                  = "Invalid priority entered."
	FileNotSentToServer              = "No file has been sent to the server."
	InvalidQuantity                  = "The entered quantity is not allowed."
	PasswordMismatch                 = "Password confirmation does not match the password."
	InvalidFileType                  = "Invalid file type"
	InvalidNationalCode              = "Invalid national code entered."
	InvalidPostalCode                = "Invalid postal code entered."
	RecoveryPasswordReqDone          = "If the provided information is accurate, an email has been dispatched to facilitate the process of password recovery for your account."
	TimeGreaterThanNow               = "The entered date time must be greater than now."
)
