package codes

import (
	"shorturl/wangjian-zero/grpc/codes"
	"shorturl/wangjian-zero/grpc/status"
)

// Acceptable checks if given error is acceptable.
func Acceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	default:
		return true
	}
}
