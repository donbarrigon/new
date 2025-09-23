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
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     any    `json:"errors,omitempty"`
}

func New(code int, message string, err any) *HttpError {
	return &HttpError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ================================
// Funciones para la interfaz de Errror
// ================================

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Errors)
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

func BadRequest(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusBadRequest,
		Message: "Solicitud incorrecta",
		Err:     err,
	}
}

func Unauthorized(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnauthorized,
		Message: "No autorizado",
		Err:     err,
	}
}

func PaymentRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusPaymentRequired,
		Message: "Pago requerido",
		Err:     err,
	}
}

func Forbidden(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusForbidden,
		Message: "Prohibido",
		Err:     err,
	}
}

func NotFound(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotFound,
		Message: "No encontrado",
		Err:     err,
	}
}

func MethodNotAllowed(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusMethodNotAllowed,
		Message: "Método no permitido",
		Err:     err,
	}
}

func NotAcceptable(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotAcceptable,
		Message: "No aceptable",
		Err:     err,
	}
}

func ProxyAuthRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusProxyAuthRequired,
		Message: "Autenticación de proxy requerida",
		Err:     err,
	}
}

func RequestTimeout(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestTimeout,
		Message: "Tiempo de solicitud agotado",
		Err:     err,
	}
}

func Conflict(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusConflict,
		Message: "Conflicto",
		Err:     err,
	}
}

func Gone(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusGone,
		Message: "Recurso no disponible",
		Err:     err,
	}
}

func LengthRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusLengthRequired,
		Message: "Longitud requerida",
		Err:     err,
	}
}

func PreconditionFailed(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusPreconditionFailed,
		Message: "Precondición fallida",
		Err:     err,
	}
}

func RequestEntityTooLarge(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestEntityTooLarge,
		Message: "Entidad de solicitud demasiado grande",
		Err:     err,
	}
}

func RequestURITooLong(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestURITooLong,
		Message: "URI demasiado larga",
		Err:     err,
	}
}

func UnsupportedMediaType(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnsupportedMediaType,
		Message: "Tipo de medio no soportado",
		Err:     err,
	}
}

func RequestedRangeNotSatisfiable(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestedRangeNotSatisfiable,
		Message: "Rango solicitado no disponible",
		Err:     err,
	}
}

func ExpectationFailed(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusExpectationFailed,
		Message: "Expectativa fallida",
		Err:     err,
	}

}

func ImATeapot(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusTeapot,
		Message: "Soy una tetera",
		Err:     err,
	}
}

func MisdirectedRequest(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusMisdirectedRequest,
		Message: "Solicitud mal dirigida",
		Err:     err,
	}
}

func UnprocessableEntity(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Entidad no procesable",
		Err:     err,
	}
}

func Locked(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusLocked,
		Message: "Recurso bloqueado",
		Err:     err,
	}
}

func FailedDependency(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusFailedDependency,
		Message: "Dependencia fallida",
		Err:     err,
	}
}

func TooEarly(err any) *HttpError {
	return &HttpError{Code: http.StatusTooEarly,
		Message: "Solicitud demasiado temprana",
		Err:     err,
	}
}

func UpgradeRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUpgradeRequired,
		Message: "Actualización requerida",
		Err:     err,
	}
}

func PreconditionRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusPreconditionRequired,
		Message: "Precondición requerida",
		Err:     err,
	}
}

func TooManyRequests(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusTooManyRequests,
		Message: "Demasiadas solicitudes",
		Err:     err,
	}
}

func RequestHeaderFieldsTooLarge(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestHeaderFieldsTooLarge,
		Message: "Encabezados de solicitud demasiado grandes",
		Err:     err,
	}
}

func UnavailableForLegalReasons(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnavailableForLegalReasons,
		Message: "No disponible por razones legales",
		Err:     err,
	}
}

// ================================
// 5xx Errores del servidor
// ================================

func Internal(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusInternalServerError,
		Message: "Error interno del servidor",
		Err:     err,
	}
}

func NotImplemented(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotImplemented,
		Message: "No implementado",
		Err:     err,
	}
}

func BadGateway(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusBadGateway,
		Message: "Puerta de enlace incorrecta",
		Err:     err,
	}
}

func ServiceUnavailable(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusServiceUnavailable,
		Message: "Servicio no disponible",
		Err:     err,
	}
}

func GatewayTimeout(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusGatewayTimeout,
		Message: "Tiempo de espera de la puerta de enlace agotado",
		Err:     err,
	}
}

func HTTPVersionNotSupported(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusHTTPVersionNotSupported,
		Message: "Versión HTTP no soportada",
		Err:     err,
	}
}

func VariantAlsoNegotiates(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusVariantAlsoNegotiates,
		Message: "Negociación de contenido fallida",
		Err:     err,
	}
}

func InsufficientStorage(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusInsufficientStorage,
		Message: "Almacenamiento insuficiente",
		Err:     err,
	}
}

func LoopDetected(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusLoopDetected,
		Message: "Bucle detectado",
		Err:     err,
	}
}

func NotExtended(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotExtended,
		Message: "Extensión requerida",
		Err:     err,
	}
}

func NetworkAuthenticationRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNetworkAuthenticationRequired,
		Message: "Autenticación de red requerida",
		Err:     err,
	}
}
