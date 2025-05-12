package errors

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"encoding/json"
)

func GrpcErrorToJson(err *pb.ErrorResponse) string {
	jsonData, _ := json.Marshal(err)
	return string(jsonData)
}
