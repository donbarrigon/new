package err

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func HexID(e any) *HttpError {
	return &HttpError{
		Status:  http.StatusBadRequest,
		Message: "El identificador no es válido",
		Err:     errorData(e),
	}
}

// Convierte los errores del driver de mongo a errores HTTP
func Mongo(e error) *HttpError {
	if e == nil {
		return nil
	}

	// No document found
	if e == mongo.ErrNoDocuments {
		return &HttpError{Status: NOT_FOUND, Message: "No encontramos lo que buscas", Err: errorData(e)}
	}

	// Client disconnected
	if e == mongo.ErrClientDisconnected {
		return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "Se perdió la conexión con la base de datos", Err: errorData(e)}
	}

	// Deadline exceeded (timeout)
	if errors.Is(e, context.DeadlineExceeded) {
		return &HttpError{Status: REQUEST_TIMEOUT, Message: "La operación tardó demasiado", Err: errorData(e)}
	}

	// Context canceled
	if errors.Is(e, context.Canceled) {
		return &HttpError{Status: BAD_REQUEST, Message: "La operación fue cancelada", Err: errorData(e)}
	}

	// Handle WriteException for more detailed write errors
	var writeException mongo.WriteException
	if errors.As(e, &writeException) {
		return handleWriteException(writeException)
	}

	// Handle BulkWriteException
	var bulkWriteException mongo.BulkWriteException
	if errors.As(e, &bulkWriteException) {
		return handleBulkWriteException(bulkWriteException)
	}

	// Handle CommandError
	var commandError mongo.CommandError
	if errors.As(e, &commandError) {
		return handleCommandError(commandError)
	}

	// Handle ServerError
	var serverError mongo.ServerError
	if errors.As(e, &serverError) {
		return handleServerError(serverError)
	}

	// Duplicate key (legacy check as fallback)
	if mongo.IsDuplicateKeyError(e) {
		return &HttpError{Status: CONFLICT, Message: "Este registro ya existe", Err: errorData(e)}
	}

	// Check for network/connection errors by error message
	errorMsg := strings.ToLower(e.Error())
	switch {
	case strings.Contains(errorMsg, "connection refused"),
		strings.Contains(errorMsg, "no reachable servers"),
		strings.Contains(errorMsg, "server selection timeout"):
		return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "No pudimos conectar con la base de datos", Err: errorData(e)}

	case strings.Contains(errorMsg, "authentication failed"):
		return &HttpError{Status: UNAUTHORIZED, Message: "Fallo la autenticación con la base de datos", Err: errorData(e)}

	case strings.Contains(errorMsg, "not authorized"):
		return &HttpError{Status: FORBIDDEN, Message: "No tienes permisos suficientes", Err: errorData(e)}

	case strings.Contains(errorMsg, "invalid namespace"):
		return &HttpError{Status: BAD_REQUEST, Message: "El nombre de la colección no es válido", Err: errorData(e)}

	case strings.Contains(errorMsg, "exceeds maximum"):
		return &HttpError{Status: BAD_REQUEST, Message: "Los datos son demasiado grandes", Err: errorData(e)}
	}

	// Default case → Internal Server Error
	return Internal(e)
}

// handleWriteException processes MongoDB write exceptions
func handleWriteException(we mongo.WriteException) *HttpError {
	for _, writeError := range we.WriteErrors {
		switch writeError.Code {
		case 11000, 11001: // Duplicate key errors
			return &HttpError{Status: CONFLICT, Message: "Este registro ya existe", Err: errorData(we)}
		case 2: // BadValue
			return &HttpError{Status: BAD_REQUEST, Message: "El valor de un campo no es válido", Err: errorData(we)}
		case 9: // FailedToParse
			return &HttpError{Status: BAD_REQUEST, Message: "No pudimos procesar los datos", Err: errorData(we)}
		case 14: // TypeMismatch
			return &HttpError{Status: BAD_REQUEST, Message: "El tipo de dato no coincide", Err: errorData(we)}
		case 16755: // Location error
			return &HttpError{Status: BAD_REQUEST, Message: "La ubicación geográfica no es válida", Err: errorData(we)}
		case 17280: // KeyTooLong
			return &HttpError{Status: BAD_REQUEST, Message: "El identificador es demasiado largo", Err: errorData(we)}
		case 10334: // BSONObjectTooLarge
			return &HttpError{Status: BAD_REQUEST, Message: "El documento es demasiado grande", Err: errorData(we)}
		}
	}

	// Check write concern errors
	if we.WriteConcernError != nil {
		switch we.WriteConcernError.Code {
		case 64: // WriteConcernFailed
			return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "No se pudo confirmar la escritura", Err: errorData(we)}
		case 79: // UnknownReplWriteConcern
			return &HttpError{Status: BAD_REQUEST, Message: "La configuración de escritura no es válida", Err: errorData(we)}
		}
	}

	return &HttpError{Status: INTERNAL, Message: "No pudimos guardar los datos", Err: errorData(we)}
}

// handleBulkWriteException processes bulk write exceptions
func handleBulkWriteException(bwe mongo.BulkWriteException) *HttpError {
	// Check for duplicate key errors in bulk operations
	for _, writeError := range bwe.WriteErrors {
		if writeError.Code == 11000 || writeError.Code == 11001 {
			return &HttpError{Status: CONFLICT, Message: "Algunos registros ya existen", Err: errorData(bwe)}
		}
	}

	// Check write concern errors
	if bwe.WriteConcernError != nil {
		return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "No se pudieron confirmar todos los cambios", Err: errorData(bwe)}
	}

	return &HttpError{Status: BAD_REQUEST, Message: "No se pudieron guardar todos los registros", Err: errorData(bwe)}
}

// handleCommandError processes MongoDB command errors
func handleCommandError(ce mongo.CommandError) *HttpError {
	switch ce.Code {
	case 2: // BadValue
		return &HttpError{Status: BAD_REQUEST, Message: "Uno de los parámetros no es válido", Err: errorData(ce)}
	case 9: // FailedToParse
		return &HttpError{Status: BAD_REQUEST, Message: "No pudimos entender la solicitud", Err: errorData(ce)}
	case 13: // Unauthorized
		return &HttpError{Status: FORBIDDEN, Message: "No tienes permisos para hacer esto", Err: errorData(ce)}
	case 18: // AuthenticationFailed
		return &HttpError{Status: UNAUTHORIZED, Message: "La autenticación falló", Err: errorData(ce)}
	case 26: // NamespaceNotFound
		return &HttpError{Status: NOT_FOUND, Message: "La colección no existe", Err: errorData(ce)}
	case 59: // CommandNotFound
		return &HttpError{Status: BAD_REQUEST, Message: "La operación no es válida", Err: errorData(ce)}
	case 61: // ShardKeyNotFound
		return &HttpError{Status: BAD_REQUEST, Message: "Falta la clave de distribución", Err: errorData(ce)}
	case 72: // InvalidOptions
		return &HttpError{Status: BAD_REQUEST, Message: "Las opciones no son válidas", Err: errorData(ce)}
	case 96: // OperationFailed
		return &HttpError{Status: INTERNAL, Message: "La operación falló", Err: errorData(ce)}
	case 11600: // InterruptedAtShutdown
		return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "El servidor se está reiniciando", Err: errorData(ce)}
	case 11601: // Interrupted
		return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "La operación fue interrumpida", Err: errorData(ce)}
	case 13435: // ShardKeyTooBig
		return &HttpError{Status: BAD_REQUEST, Message: "La clave de distribución es demasiado grande", Err: errorData(ce)}
	case 16550: // DocumentValidationFailure
		return &HttpError{Status: BAD_REQUEST, Message: "Los datos no cumplen con las reglas de validación", Err: errorData(ce)}
	case 50: // MaxTimeMSExpired
		return &HttpError{Status: REQUEST_TIMEOUT, Message: "La operación tomó demasiado tiempo", Err: errorData(ce)}
	}

	return &HttpError{Status: INTERNAL, Message: "La operación falló", Err: errorData(ce)}
}

// handleServerError processes general MongoDB server errors
func handleServerError(se mongo.ServerError) *HttpError {
	// Server errors are typically infrastructure issues
	return &HttpError{Status: SERVICE_UNAVAILABLE, Message: "Hay un problema con el servidor de la base de datos", Err: errorData(se)}
}
