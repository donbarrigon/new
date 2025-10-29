package err

import (
	"donbarrigon/new/internal/utils/config"
	"net/http"
	"reflect"
)

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
// Funciones para la interfaz de error
// ================================

func (e *HttpError) Error() string {
	return e.Message
}

// ================================
// 4xx Errores del cliente
// ================================
func BadRequest(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusBadRequest,
		Message: "Algo no está bien con tu solicitud",
		Err:     errorData(e),
	}
}

func Unauthorized(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnauthorized,
		Message: "No tiene autorizacion para hacer esto",
		Err:     errorData(e),
	}
}

func PaymentRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPaymentRequired,
		Message: "Se requiere un pago para acceder",
		Err:     errorData(e),
	}
}

func Forbidden(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusForbidden,
		Message: "No tienes permiso para hacer esto",
		Err:     errorData(e),
	}
}

func NotFound(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotFound,
		Message: "No encontramos lo que buscas",
		Err:     errorData(e),
	}
}

func MethodNotAllowed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusMethodNotAllowed,
		Message: "Esta acción no está permitida",
		Err:     errorData(e),
	}
}

func NotAcceptable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotAcceptable,
		Message: "No podemos procesar tu solicitud en este formato",
		Err:     errorData(e),
	}
}

func ProxyAuthRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusProxyAuthRequired,
		Message: "Se requiere autenticación del proxy",
		Err:     errorData(e),
	}
}

func RequestTimeout(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestTimeout,
		Message: "La solicitud tardó demasiado tiempo",
		Err:     errorData(e),
	}
}

func Conflict(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusConflict,
		Message: "Hay un conflicto con los datos existentes",
		Err:     errorData(e),
	}
}

func Gone(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusGone,
		Message: "Este recurso ya no está disponible",
		Err:     errorData(e),
	}
}

func LengthRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLengthRequired,
		Message: "Falta información del tamaño de la solicitud",
		Err:     errorData(e),
	}
}

func PreconditionFailed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPreconditionFailed,
		Message: "No se cumplen las condiciones necesarias",
		Err:     errorData(e),
	}
}

func RequestEntityTooLarge(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestEntityTooLarge,
		Message: "El archivo o datos son demasiado grandes",
		Err:     errorData(e),
	}
}

func RequestURITooLong(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestURITooLong,
		Message: "La dirección es demasiado larga",
		Err:     errorData(e),
	}
}

func UnsupportedMediaType(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnsupportedMediaType,
		Message: "El tipo de archivo no es compatible",
		Err:     errorData(e),
	}
}

func RequestedRangeNotSatisfiable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestedRangeNotSatisfiable,
		Message: "El rango solicitado no es válido",
		Err:     errorData(e),
	}
}

func ExpectationFailed(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusExpectationFailed,
		Message: "No se pudo cumplir con lo esperado",
		Err:     errorData(e),
	}
}

func ImATeapot(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTeapot,
		Message: "Soy una tetera, no puedo hacer café",
		Err:     errorData(e),
	}
}

func MisdirectedRequest(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusMisdirectedRequest,
		Message: "Tu solicitud fue enviada al lugar equivocado",
		Err:     errorData(e),
	}
}

func UnprocessableEntity(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnprocessableEntity,
		Message: "No pudimos procesar la información que enviaste",
		Err:     errorData(e),
	}
}

func Locked(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLocked,
		Message: "Este recurso está bloqueado",
		Err:     errorData(e),
	}
}

func FailedDependency(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusFailedDependency,
		Message: "Algo que necesitamos no funcionó correctamente",
		Err:     errorData(e),
	}
}

func TooEarly(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTooEarly,
		Message: "Es demasiado pronto para esta solicitud",
		Err:     errorData(e),
	}
}

func UpgradeRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUpgradeRequired,
		Message: "Necesitas actualizar tu aplicación",
		Err:     errorData(e),
	}
}

func PreconditionRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusPreconditionRequired,
		Message: "Faltan algunos requisitos previos",
		Err:     errorData(e),
	}
}

func TooManyRequests(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusTooManyRequests,
		Message: "Has hecho demasiadas solicitudes, espera un momento",
		Err:     errorData(e),
	}
}

func RequestHeaderFieldsTooLarge(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusRequestHeaderFieldsTooLarge,
		Message: "La información de tu solicitud es demasiado grande",
		Err:     errorData(e),
	}
}

func UnavailableForLegalReasons(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusUnavailableForLegalReasons,
		Message: "No disponible por razones legales",
		Err:     errorData(e),
	}
}

// ================================
// 5xx Server errors
// ================================

func Internal(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusInternalServerError,
		Message: "Algo salió mal de nuestro lado",
		Err:     errorData(e),
	}
}

func Panic(e any, stack string) *HttpError {
	return &HttpError{
		Status:  http.StatusInternalServerError,
		Message: "Ups, algo salió muy mal",
		Err:     map[string]any{"error": errorData(e), "stack": stack},
	}
}

func NotImplemented(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotImplemented,
		Message: "Esta función aún no está disponible",
		Err:     errorData(e),
	}
}

func BadGateway(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusBadGateway,
		Message: "Hubo un problema con nuestros servidores",
		Err:     errorData(e),
	}
}

func ServiceUnavailable(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusServiceUnavailable,
		Message: "El servicio no está disponible en este momento",
		Err:     errorData(e),
	}
}

func GatewayTimeout(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusGatewayTimeout,
		Message: "El servidor tardó demasiado en responder",
		Err:     errorData(e),
	}
}

func HTTPVersionNotSupported(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusHTTPVersionNotSupported,
		Message: "La versión de tu navegador no es compatible",
		Err:     errorData(e),
	}
}

func VariantAlsoNegotiates(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusVariantAlsoNegotiates,
		Message: "Hubo un error en la negociación del servidor",
		Err:     errorData(e),
	}
}

func InsufficientStorage(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusInsufficientStorage,
		Message: "No hay suficiente espacio disponible",
		Err:     errorData(e),
	}
}

func LoopDetected(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusLoopDetected,
		Message: "Se detectó un bucle infinito en el servidor",
		Err:     errorData(e),
	}
}

func NotExtended(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNotExtended,
		Message: "Faltan extensiones necesarias en el servidor",
		Err:     errorData(e),
	}
}

func NetworkAuthenticationRequired(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusNetworkAuthenticationRequired,
		Message: "Necesitas autenticarte en la red",
		Err:     errorData(e),
	}
}

func errorData(e any) any {
	if e == nil {
		return nil
	}

	if !config.AppDebug {
		return nil
	}

	v := reflect.ValueOf(e)
	if v.Kind() == reflect.Pointer {
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
