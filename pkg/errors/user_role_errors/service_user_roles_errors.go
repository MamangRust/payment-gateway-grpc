package userrole_errors

import (
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"net/http"
)

var (
	ErrFailedAssignRoleToUser = errors.NewErrorResponse("Failed to assign role to user", http.StatusInternalServerError)
	ErrFailedRemoveRole       = errors.NewErrorResponse("Failed to remove role from user", http.StatusInternalServerError)
)
