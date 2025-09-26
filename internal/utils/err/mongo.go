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
		Message: "Invalid hex id",
		Err:     e,
	}
}

// Mongo maps MongoDB driver errors to HttpError
func Mongo(e error) *HttpError {
	if e == nil {
		return nil
	}

	// No document found
	if e == mongo.ErrNoDocuments {
		return NotFound("document not found")
	}

	// Client disconnected
	if e == mongo.ErrClientDisconnected {
		return ServiceUnavailable("database client disconnected")
	}

	// Deadline exceeded (timeout)
	if errors.Is(e, context.DeadlineExceeded) {
		return RequestTimeout("operation timed out")
	}

	// Context canceled
	if errors.Is(e, context.Canceled) {
		return BadRequest("operation canceled")
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
		return Conflict("duplicate key error")
	}

	// Check for network/connection errors by error message
	errorMsg := strings.ToLower(e.Error())
	switch {
	case strings.Contains(errorMsg, "connection refused"),
		strings.Contains(errorMsg, "no reachable servers"),
		strings.Contains(errorMsg, "server selection timeout"):
		return ServiceUnavailable("database connection failed")

	case strings.Contains(errorMsg, "authentication failed"):
		return Unauthorized("database authentication failed")

	case strings.Contains(errorMsg, "not authorized"):
		return Forbidden("insufficient database privileges")

	case strings.Contains(errorMsg, "invalid namespace"):
		return BadRequest("invalid database or collection name")

	case strings.Contains(errorMsg, "exceeds maximum"):
		return BadRequest("request exceeds maximum allowed size")
	}

	// Default case → Internal Server Error
	return Internal(e.Error())
}

// handleWriteException processes MongoDB write exceptions
func handleWriteException(we mongo.WriteException) *HttpError {
	for _, writeError := range we.WriteErrors {
		switch writeError.Code {
		case 11000, 11001: // Duplicate key errors
			return Conflict("duplicate key violation")
		case 2: // BadValue
			return BadRequest("invalid field value")
		case 9: // FailedToParse
			return BadRequest("failed to parse document")
		case 14: // TypeMismatch
			return BadRequest("field type mismatch")
		case 16755: // Location error
			return BadRequest("invalid geospatial data")
		case 17280: // KeyTooLong
			return BadRequest("key too long")
		case 10334: // BSONObjectTooLarge
			return BadRequest("document too large")
		}
	}

	// Check write concern errors
	if we.WriteConcernError != nil {
		switch we.WriteConcernError.Code {
		case 64: // WriteConcernFailed
			return ServiceUnavailable("write concern failed")
		case 79: // UnknownReplWriteConcern
			return BadRequest("unknown write concern")
		}
	}

	return Internal("write operation failed: " + we.Error())
}

// handleBulkWriteException processes bulk write exceptions
func handleBulkWriteException(bwe mongo.BulkWriteException) *HttpError {
	// Check for duplicate key errors in bulk operations
	for _, writeError := range bwe.WriteErrors {
		if writeError.Code == 11000 || writeError.Code == 11001 {
			return Conflict("bulk operation contains duplicate keys")
		}
	}

	// Check write concern errors
	if bwe.WriteConcernError != nil {
		return ServiceUnavailable("bulk write concern failed")
	}

	return BadRequest("bulk write operation failed")
}

// handleCommandError processes MongoDB command errors
func handleCommandError(ce mongo.CommandError) *HttpError {
	switch ce.Code {
	case 2: // BadValue
		return BadRequest("invalid command parameter")
	case 9: // FailedToParse
		return BadRequest("failed to parse command")
	case 13: // Unauthorized
		return Forbidden("insufficient privileges for command")
	case 18: // AuthenticationFailed
		return Unauthorized("authentication failed")
	case 26: // NamespaceNotFound
		return NotFound("database or collection not found")
	case 59: // CommandNotFound
		return BadRequest("unknown command")
	case 61: // ShardKeyNotFound
		return BadRequest("shard key not found")
	case 72: // InvalidOptions
		return BadRequest("invalid command options")
	case 96: // OperationFailed
		return Internal("database operation failed")
	case 11600: // InterruptedAtShutdown
		return ServiceUnavailable("database shutting down")
	case 11601: // Interrupted
		return ServiceUnavailable("operation interrupted")
	case 13435: // ShardKeyTooBig
		return BadRequest("shard key too large")
	case 16550: // DocumentValidationFailure
		return BadRequest("document validation failed")
	case 50: // MaxTimeMSExpired
		return RequestTimeout("operation exceeded time limit")
	}

	return Internal("command failed: " + ce.Error())
}

// handleServerError processes general MongoDB server errors
func handleServerError(se mongo.ServerError) *HttpError {
	// Server errors are typically infrastructure issues
	return ServiceUnavailable("database server error: " + se.Error())
}
