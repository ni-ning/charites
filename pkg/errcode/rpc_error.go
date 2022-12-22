package errcode

import (
	pb "charites/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// func ToRPCError(err *Error) error {
// 	s := status.New(ToRPCCode(err.Code()), err.Msg())
// 	return s.Err()
// }

func ToRPCError(err *Error) error {
	// proto定义的只不过是业务错误码响应格式
	pbErr := &pb.Error{Code: int32(err.Code()), Message: err.Msg()}
	// status.New中的Status有Details这个字段所以才能正常返回
	s, _ := status.New(ToRPCCode(err.Code()), err.Msg()).WithDetails(pbErr)
	return s.Err()
}

func ToRPCCode(code int) codes.Code {
	var statusCode codes.Code
	switch code {
	case Fail.Code():
		statusCode = codes.Internal
	case InvalidParams.Code():
		statusCode = codes.InvalidArgument
	case Unauthorized.Code():
		statusCode = codes.Unauthenticated
	case AccessDenied.Code():
		statusCode = codes.PermissionDenied
	case DeadlineExceeded.Code():
		statusCode = codes.DeadlineExceeded
	case NotFound.Code():
		statusCode = codes.NotFound
	case LimitExceed.Code():
		statusCode = codes.ResourceExhausted
	case MethodNotAllowed.Code():
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}
	return statusCode
}

func ToRPCStatus(code int, msg string) *status.Status {
	pbErr := &pb.Error{Code: int32(code), Message: msg}
	s, _ := status.New(ToRPCCode(code), msg).WithDetails(pbErr)
	return s
}

func FromError(err error) *status.Status {
	s, _ := status.FromError(err)
	return s
}
