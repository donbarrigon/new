package err

import "net/http"

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (e *HttpError) Error() string {
	return e.Message
}

// =========================
// 4xx Errores del cliente
// =========================

func BadRequest(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusBadRequest,
		Message: "Solicitud incorrecta",
		Errors:  err,
	}
}

func Unauthorized(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnauthorized,
		Message: "No autorizado",
		Errors:  err,
	}
}

func PaymentRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusPaymentRequired,
		Message: "Pago requerido",
		Errors:  err,
	}
}

func Forbidden(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusForbidden,
		Message: "Prohibido",
		Errors:  err,
	}
}

func NotFound(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotFound,
		Message: "No encontrado",
		Errors:  err,
	}
}

func MethodNotAllowed(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusMethodNotAllowed,
		Message: "Método no permitido",
		Errors:  err,
	}
}

func NotAcceptable(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotAcceptable,
		Message: "No aceptable",
		Errors:  err,
	}
}

func ProxyAuthRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusProxyAuthRequired,
		Message: "Autenticación de proxy requerida",
		Errors:  err,
	}
}

func RequestTimeout(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestTimeout,
		Message: "Tiempo de solicitud agotado",
		Errors:  err,
	}
}

func Conflict(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusConflict,
		Message: "Conflicto",
		Errors:  err,
	}
}

func Gone(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusGone,
		Message: "Recurso no disponible",
		Errors:  err,
	}
}

func LengthRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusLengthRequired,
		Message: "Longitud requerida",
		Errors:  err,
	}
}

func PreconditionFailed(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusPreconditionFailed,
		Message: "Precondición fallida",
		Errors:  err,
	}
}

func RequestEntityTooLarge(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestEntityTooLarge,
		Message: "Entidad de solicitud demasiado grande",
		Errors:  err,
	}
}

func RequestURITooLong(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestURITooLong,
		Message: "URI demasiado larga",
		Errors:  err,
	}
}

func UnsupportedMediaType(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnsupportedMediaType,
		Message: "Tipo de medio no soportado",
		Errors:  err,
	}
}

func RequestedRangeNotSatisfiable(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestedRangeNotSatisfiable,
		Message: "Rango solicitado no disponible",
		Errors:  err,
	}
}

func ExpectationFailed(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusExpectationFailed,
		Message: "Expectativa fallida",
		Errors:  err,
	}

}

func ImATeapot(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusTeapot,
		Message: "Soy una tetera",
		Errors:  err,
	}
}

func MisdirectedRequest(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusMisdirectedRequest,
		Message: "Solicitud mal dirigida",
		Errors:  err,
	}
}

func UnprocessableEntity(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Entidad no procesable",
		Errors:  err,
	}
}

func Locked(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusLocked,
		Message: "Recurso bloqueado",
		Errors:  err,
	}
}

func FailedDependency(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusFailedDependency,
		Message: "Dependencia fallida",
		Errors:  err,
	}
}

func TooEarly(err any) *HttpError {
	return &HttpError{Code: http.StatusTooEarly,
		Message: "Solicitud demasiado temprana",
		Errors:  err,
	}
}

func UpgradeRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUpgradeRequired,
		Message: "Actualización requerida",
		Errors:  err,
	}
}

func PreconditionRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusPreconditionRequired,
		Message: "Precondición requerida",
		Errors:  err,
	}
}

func TooManyRequests(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusTooManyRequests,
		Message: "Demasiadas solicitudes",
		Errors:  err,
	}
}

func RequestHeaderFieldsTooLarge(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusRequestHeaderFieldsTooLarge,
		Message: "Encabezados de solicitud demasiado grandes",
		Errors:  err,
	}
}

func UnavailableForLegalReasons(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusUnavailableForLegalReasons,
		Message: "No disponible por razones legales",
		Errors:  err,
	}
}

// =========================
// 5xx Errores del servidor
// =========================

func InternalServerError(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusInternalServerError,
		Message: "Error interno del servidor",
		Errors:  err,
	}
}

func NotImplemented(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotImplemented,
		Message: "No implementado",
		Errors:  err,
	}
}

func BadGateway(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusBadGateway,
		Message: "Puerta de enlace incorrecta",
		Errors:  err,
	}
}

func ServiceUnavailable(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusServiceUnavailable,
		Message: "Servicio no disponible",
		Errors:  err,
	}
}

func GatewayTimeout(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusGatewayTimeout,
		Message: "Tiempo de espera de la puerta de enlace agotado",
		Errors:  err,
	}
}

func HTTPVersionNotSupported(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusHTTPVersionNotSupported,
		Message: "Versión HTTP no soportada",
		Errors:  err,
	}
}

func VariantAlsoNegotiates(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusVariantAlsoNegotiates,
		Message: "Negociación de contenido fallida",
		Errors:  err,
	}
}

func InsufficientStorage(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusInsufficientStorage,
		Message: "Almacenamiento insuficiente",
		Errors:  err,
	}
}

func LoopDetected(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusLoopDetected,
		Message: "Bucle detectado",
		Errors:  err,
	}
}

func NotExtended(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNotExtended,
		Message: "Extensión requerida",
		Errors:  err,
	}
}

func NetworkAuthenticationRequired(err any) *HttpError {
	return &HttpError{
		Code:    http.StatusNetworkAuthenticationRequired,
		Message: "Autenticación de red requerida",
		Errors:  err,
	}
}
