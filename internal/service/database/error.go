package database

import (
    "fmt"
    "strings"
    "net/http"
    "database/sql"

    "github.com/jsirianni/systemstat/internal/log"

    "github.com/pkg/errors"
)

func errorToHTTPStatus(err error) int32 {
    if isErrNoRows(err) {
        return http.StatusNotFound
    }

    if isErrInvalidUUID(err) {
        return http.StatusBadRequest
    }

    log.Error(errors.New(fmt.Sprintf("failed to convert sql error to http status, unknown sql error: %s", err.Error())))
    return http.StatusInternalServerError
}

func isErrInvalidUUID(err error) bool {
    const substr = "invalid input syntax for type uuid"
    return strings.Contains(err.Error(), substr)
}

func isErrNoRows(err error) bool {
    return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}
