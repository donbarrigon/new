package err

import (
	"donbarrigon/new/internal/utils/lang"
	"fmt"
	"net/http"
)

type Error interface {
	error
	Errors(l string) *HttpError
}

type HttpError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Err     any    `json:"errors,omitempty"`
}

func New(status int, message string, e any) *HttpError {
	return &HttpError{
		Status:  status,
		Message: message,
		Err:     e,
	}
}

// ================================
// Funciones para la interfaz de Errror
// ================================

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err)
}

func (e *HttpError) Errors(l string) *HttpError {
	e.Message = lang.T(l, e.Message, nil)
	if s, ok := e.Err.(string); ok {
		e.Err = lang.T(l, s, nil)
	}
	return e
}

// ================================
// 4xx Errores del cliente
// ================================
func BadRequest(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusBadRequest,
		Message: "Bad request",
		Err:     e,
	}
}

func Unauthorized(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Err:     e,
	}
}

func PaymentRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPaymentRequired,
		Message: "Payment required",
		Err:     e,
	}
}

func Forbidden(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusForbidden,
		Message: "Forbidden",
		Err:     e,
	}
}

func NotFound(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotFound,
		Message: "Not found",
		Err:     e,
	}
}

func MethodNotAllowed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method not allowed",
		Err:     e,
	}
}

func NotAcceptable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotAcceptable,
		Message: "Not acceptable",
		Err:     e,
	}
}

func ProxyAuthRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusProxyAuthRequired,
		Message: "Proxy authentication required",
		Err:     e,
	}
}

func RequestTimeout(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestTimeout,
		Message: "Request timeout",
		Err:     e,
	}
}

func Conflict(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusConflict,
		Message: "Conflict",
		Err:     e,
	}
}

func Gone(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusGone,
		Message: "Gone",
		Err:     e,
	}
}

func LengthRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLengthRequired,
		Message: "Length required",
		Err:     e,
	}
}

func PreconditionFailed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPreconditionFailed,
		Message: "Precondition failed",
		Err:     e,
	}
}

func RequestEntityTooLarge(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestEntityTooLarge,
		Message: "Payload too large",
		Err:     e,
	}
}

func RequestURITooLong(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestURITooLong,
		Message: "URI too long",
		Err:     e,
	}
}

func UnsupportedMediaType(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnsupportedMediaType,
		Message: "Unsupported media type",
		Err:     e,
	}
}

func RequestedRangeNotSatisfiable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestedRangeNotSatisfiable,
		Message: "Range not satisfiable",
		Err:     e,
	}
}

func ExpectationFailed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusExpectationFailed,
		Message: "Expectation failed",
		Err:     e,
	}
}

func ImATeapot(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTeapot,
		Message: "I'm a teapot",
		Err:     e,
	}
}

func MisdirectedRequest(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusMisdirectedRequest,
		Message: "Misdirected request",
		Err:     e,
	}
}

func UnprocessableEntity(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Unprocessable entity",
		Err:     e,
	}
}

func Locked(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLocked,
		Message: "Locked",
		Err:     e,
	}
}

func FailedDependency(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusFailedDependency,
		Message: "Failed dependency",
		Err:     e,
	}
}

func TooEarly(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTooEarly,
		Message: "Too early",
		Err:     e,
	}
}

func UpgradeRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUpgradeRequired,
		Message: "Upgrade required",
		Err:     e,
	}
}

func PreconditionRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPreconditionRequired,
		Message: "Precondition required",
		Err:     e,
	}
}

func TooManyRequests(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTooManyRequests,
		Message: "Too many requests",
		Err:     e,
	}
}

func RequestHeaderFieldsTooLarge(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestHeaderFieldsTooLarge,
		Message: "Request header fields too large",
		Err:     e,
	}
}

func UnavailableForLegalReasons(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnavailableForLegalReasons,
		Message: "Unavailable for legal reasons",
		Err:     e,
	}
}

// ================================
// 5xx Server errors
// ================================

func Internal(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
		Err:     e,
	}
}

func NotImplemented(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotImplemented,
		Message: "Not implemented",
		Err:     e,
	}
}

func BadGateway(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusBadGateway,
		Message: "Bad gateway",
		Err:     e,
	}
}

func ServiceUnavailable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusServiceUnavailable,
		Message: "Service unavailable",
		Err:     e,
	}
}

func GatewayTimeout(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusGatewayTimeout,
		Message: "Gateway timeout",
		Err:     e,
	}
}

func HTTPVersionNotSupported(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusHTTPVersionNotSupported,
		Message: "HTTP version not supported",
		Err:     e,
	}
}

func VariantAlsoNegotiates(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusVariantAlsoNegotiates,
		Message: "Variant also negotiates",
		Err:     e,
	}
}

func InsufficientStorage(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusInsufficientStorage,
		Message: "Insufficient storage",
		Err:     e,
	}
}

func LoopDetected(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLoopDetected,
		Message: "Loop detected",
		Err:     e,
	}
}

func NotExtended(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotExtended,
		Message: "Not extended",
		Err:     e,
	}
}

func NetworkAuthenticationRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNetworkAuthenticationRequired,
		Message: "Network authentication required",
		Err:     e,
	}
}
