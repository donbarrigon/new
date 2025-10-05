package err

import (
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/lang"
	"fmt"
	"net/http"
	"reflect"
)

type Error interface {
	error
	Errors(l string) *HttpError
}

type HttpError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Err     any    `json:"error,omitempty"`
}

func New(status int, message string, e any) *HttpError {
	return &HttpError{
		Status:  status,
		Message: message,
		Err:     errorData(e),
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
	if config.AppDebug {
		e.Err = nil
	}
	if e.Err == nil {
		return e
	}
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
		Err:     errorData(e),
	}
}

func Unauthorized(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Err:     errorData(e),
	}
}

func PaymentRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPaymentRequired,
		Message: "Payment required",
		Err:     errorData(e),
	}
}

func Forbidden(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusForbidden,
		Message: "Forbidden",
		Err:     errorData(e),
	}
}

func NotFound(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotFound,
		Message: "Not found",
		Err:     errorData(e),
	}
}

func MethodNotAllowed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method not allowed",
		Err:     errorData(e),
	}
}

func NotAcceptable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotAcceptable,
		Message: "Not acceptable",
		Err:     errorData(e),
	}
}

func ProxyAuthRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusProxyAuthRequired,
		Message: "Proxy authentication required",
		Err:     errorData(e),
	}
}

func RequestTimeout(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestTimeout,
		Message: "Request timeout",
		Err:     errorData(e),
	}
}

func Conflict(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusConflict,
		Message: "Conflict",
		Err:     errorData(e),
	}
}

func Gone(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusGone,
		Message: "Gone",
		Err:     errorData(e),
	}
}

func LengthRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLengthRequired,
		Message: "Length required",
		Err:     errorData(e),
	}
}

func PreconditionFailed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPreconditionFailed,
		Message: "Precondition failed",
		Err:     errorData(e),
	}
}

func RequestEntityTooLarge(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestEntityTooLarge,
		Message: "Payload too large",
		Err:     errorData(e),
	}
}

func RequestURITooLong(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestURITooLong,
		Message: "URI too long",
		Err:     errorData(e),
	}
}

func UnsupportedMediaType(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnsupportedMediaType,
		Message: "Unsupported media type",
		Err:     errorData(e),
	}
}

func RequestedRangeNotSatisfiable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestedRangeNotSatisfiable,
		Message: "Range not satisfiable",
		Err:     errorData(e),
	}
}

func ExpectationFailed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusExpectationFailed,
		Message: "Expectation failed",
		Err:     errorData(e),
	}
}

func ImATeapot(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTeapot,
		Message: "I'm a teapot",
		Err:     errorData(e),
	}
}

func MisdirectedRequest(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusMisdirectedRequest,
		Message: "Misdirected request",
		Err:     errorData(e),
	}
}

func UnprocessableEntity(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Unprocessable entity",
		Err:     errorData(e),
	}
}

func Locked(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLocked,
		Message: "Locked",
		Err:     errorData(e),
	}
}

func FailedDependency(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusFailedDependency,
		Message: "Failed dependency",
		Err:     errorData(e),
	}
}

func TooEarly(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTooEarly,
		Message: "Too early",
		Err:     errorData(e),
	}
}

func UpgradeRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUpgradeRequired,
		Message: "Upgrade required",
		Err:     errorData(e),
	}
}

func PreconditionRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPreconditionRequired,
		Message: "Precondition required",
		Err:     errorData(e),
	}
}

func TooManyRequests(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTooManyRequests,
		Message: "Too many requests",
		Err:     errorData(e),
	}
}

func RequestHeaderFieldsTooLarge(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestHeaderFieldsTooLarge,
		Message: "Request header fields too large",
		Err:     errorData(e),
	}
}

func UnavailableForLegalReasons(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnavailableForLegalReasons,
		Message: "Unavailable for legal reasons",
		Err:     errorData(e),
	}
}

// ================================
// 5xx Server errors
// ================================

func Internal(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
		Err:     errorData(e),
	}
}

func Panic(e any, stack string) *HttpError {
	return &HttpError{
		Status:  http.StatusInternalServerError,
		Message: "Panic",
		Err:     map[string]any{"error": errorData(e), "stack": stack},
	}
}

func NotImplemented(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotImplemented,
		Message: "Not implemented",
		Err:     errorData(e),
	}
}

func BadGateway(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusBadGateway,
		Message: "Bad gateway",
		Err:     errorData(e),
	}
}

func ServiceUnavailable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusServiceUnavailable,
		Message: "Service unavailable",
		Err:     errorData(e),
	}
}

func GatewayTimeout(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusGatewayTimeout,
		Message: "Gateway timeout",
		Err:     errorData(e),
	}
}

func HTTPVersionNotSupported(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusHTTPVersionNotSupported,
		Message: "HTTP version not supported",
		Err:     errorData(e),
	}
}

func VariantAlsoNegotiates(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusVariantAlsoNegotiates,
		Message: "Variant also negotiates",
		Err:     errorData(e),
	}
}

func InsufficientStorage(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusInsufficientStorage,
		Message: "Insufficient storage",
		Err:     errorData(e),
	}
}

func LoopDetected(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLoopDetected,
		Message: "Loop detected",
		Err:     errorData(e),
	}
}

func NotExtended(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotExtended,
		Message: "Not extended",
		Err:     errorData(e),
	}
}

func NetworkAuthenticationRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNetworkAuthenticationRequired,
		Message: "Network authentication required",
		Err:     errorData(e),
	}
}

func errorData(e any) any {
	if e == nil {
		return nil
	}

	v := reflect.ValueOf(e)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		result := make(map[string]any)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)

			if field.IsExported() {
				fieldValue := v.Field(i)
				result[field.Name] = fieldValue.Interface()
			}
		}

		if len(result) > 0 {
			return result
		}
	}

	if er, ok := e.(error); ok {
		return er.Error()
	}

	return e
}
