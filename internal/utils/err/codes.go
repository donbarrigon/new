package err

// Códigos de error HTTP 4xx y 5xx
const (
	// --- 4xx: Errores del cliente ---
	BAD_REQUEST                     = 400
	UNAUTHORIZED                    = 401
	PAYMENT_REQUIRED                = 402
	FORBIDDEN                       = 403
	NOT_FOUND                       = 404
	METHOD_NOT_ALLOWED              = 405
	NOT_ACCEPTABLE                  = 406
	PROXY_AUTHENTICATION_REQUIRED   = 407
	REQUEST_TIMEOUT                 = 408
	CONFLICT                        = 409
	GONE                            = 410
	LENGTH_REQUIRED                 = 411
	PRECONDITION_FAILED             = 412
	PAYLOAD_TOO_LARGE               = 413
	URI_TOO_LONG                    = 414
	UNSUPPORTED_MEDIA_TYPE          = 415
	RANGE_NOT_SATISFIABLE           = 416
	EXPECTATION_FAILED              = 417
	IM_A_TEAPOT                     = 418
	MISDIRECTED_REQUEST             = 421
	UNPROCESSABLE_ENTITY            = 422
	LOCKED                          = 423
	FAILED_DEPENDENCY               = 424
	TOO_EARLY                       = 425
	UPGRADE_REQUIRED                = 426
	PRECONDITION_REQUIRED           = 428
	TOO_MANY_REQUESTS               = 429
	REQUEST_HEADER_FIELDS_TOO_LARGE = 431
	UNAVAILABLE_FOR_LEGAL_REASONS   = 451

	// --- 5xx: Errores del servidor ---
	INTERNAL                        = 500
	NOT_IMPLEMENTED                 = 501
	BAD_GATEWAY                     = 502
	SERVICE_UNAVAILABLE             = 503
	GATEWAY_TIMEOUT                 = 504
	HTTP_VERSION_NOT_SUPPORTED      = 505
	VARIANT_ALSO_NEGOTIATES         = 506
	INSUFFICIENT_STORAGE            = 507
	LOOP_DETECTED                   = 508
	NOT_EXTENDED                    = 510
	NETWORK_AUTHENTICATION_REQUIRED = 511
)

var StatusMap = map[int]string{
	// --- 4xx: Errores del cliente ---
	BAD_REQUEST:                     "Algo no está bien con tu solicitud",
	UNAUTHORIZED:                    "Necesitas iniciar sesión para continuar",
	PAYMENT_REQUIRED:                "Se requiere un pago para acceder",
	FORBIDDEN:                       "No tienes permiso para hacer esto",
	NOT_FOUND:                       "No encontramos lo que buscas",
	METHOD_NOT_ALLOWED:              "Esta acción no está permitida",
	NOT_ACCEPTABLE:                  "No podemos procesar tu solicitud en este formato",
	PROXY_AUTHENTICATION_REQUIRED:   "Se requiere autenticación del proxy",
	REQUEST_TIMEOUT:                 "La solicitud tardó demasiado tiempo",
	CONFLICT:                        "Hay un conflicto con los datos existentes",
	GONE:                            "Este recurso ya no está disponible",
	LENGTH_REQUIRED:                 "Falta información del tamaño de la solicitud",
	PRECONDITION_FAILED:             "No se cumplen las condiciones necesarias",
	PAYLOAD_TOO_LARGE:               "El archivo o datos son demasiado grandes",
	URI_TOO_LONG:                    "La dirección es demasiado larga",
	UNSUPPORTED_MEDIA_TYPE:          "El tipo de archivo no es compatible",
	RANGE_NOT_SATISFIABLE:           "El rango solicitado no es válido",
	EXPECTATION_FAILED:              "No se pudo cumplir con lo esperado",
	IM_A_TEAPOT:                     "Soy una tetera, no puedo hacer café",
	MISDIRECTED_REQUEST:             "Tu solicitud fue enviada al lugar equivocado",
	UNPROCESSABLE_ENTITY:            "No pudimos procesar la información que enviaste",
	LOCKED:                          "Este recurso está bloqueado",
	FAILED_DEPENDENCY:               "Algo que necesitamos no funcionó correctamente",
	TOO_EARLY:                       "Es demasiado pronto para esta solicitud",
	UPGRADE_REQUIRED:                "Necesitas actualizar tu aplicación",
	PRECONDITION_REQUIRED:           "Faltan algunos requisitos previos",
	TOO_MANY_REQUESTS:               "Has hecho demasiadas solicitudes, espera un momento",
	REQUEST_HEADER_FIELDS_TOO_LARGE: "La información de tu solicitud es demasiado grande",
	UNAVAILABLE_FOR_LEGAL_REASONS:   "No disponible por razones legales",

	// --- 5xx: Errores del servidor ---
	INTERNAL:                        "Algo salió mal de nuestro lado",
	NOT_IMPLEMENTED:                 "Esta función aún no está disponible",
	BAD_GATEWAY:                     "Hubo un problema con nuestros servidores",
	SERVICE_UNAVAILABLE:             "El servicio no está disponible en este momento",
	GATEWAY_TIMEOUT:                 "El servidor tardó demasiado en responder",
	HTTP_VERSION_NOT_SUPPORTED:      "La versión de tu navegador no es compatible",
	VARIANT_ALSO_NEGOTIATES:         "Hubo un error en la negociación del servidor",
	INSUFFICIENT_STORAGE:            "No hay suficiente espacio disponible",
	LOOP_DETECTED:                   "Se detectó un bucle infinito en el servidor",
	NOT_EXTENDED:                    "Faltan extensiones necesarias en el servidor",
	NETWORK_AUTHENTICATION_REQUIRED: "Necesitas autenticarte en la red",
}
