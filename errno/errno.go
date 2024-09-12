package errno

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var _ Error = (*errno)(nil)

type (
	Error interface {
		GetCode() int
		GetHttpStatusCode() int
		GetErrMsg() string

		GetData() interface{}
		GetRawData() []byte

		// Format 补充格式化输出错误信息
		Format(v ...interface{}) Error

		// WithData 设置成功时返回的数据
		WithData(data interface{}) Error

		// WithReason 错误详细描述
		WithReason(reason interface{}) Error

		WithGrpcError(err error) Error
		// WithID 设置当前请求的唯一ID
		WithID(id string) Error

		// 设置http status code
		WithHttpStatusCode(httpStatusCode int) Error

		ToString() string

		ToBytes() []byte

		ToDataBytes() []byte

		ToDataString() string

		Error() string
		GetErrorCode() string
		GetNowTime() int64
		GetRequestId() string
		Render
	}

	errno struct {
		HttpStatusCode int         `json:"-"`
		Code           int         `json:"-"` // 业务编码
		ErrorCode      string      `json:"error_code,omitempty"`
		Msg            string      `json:"message,omitempty"`    // 错误描述
		Data           interface{} `json:"data,omitempty"`       // 成功时返回的数据
		Reason         interface{} `json:"-"`                    // 错误详细描述
		RequestID      string      `json:"request_id,omitempty"` // 当前请求的唯一ID，便于问题定位
		NowTime        int64       `json:"now_time"`             // 时间戳
	}
)

func NewError(code int, errCode, errMsg string) Error {
	return errno{
		HttpStatusCode: http.StatusOK,
		ErrorCode:      errCode,
		Code:           code,
		Msg:            errMsg,
		//Data:           make(map[string]interface{}),
		NowTime: time.Now().Unix(),
	}
}

func (e errno) GetErrorCode() string {
	return e.ErrorCode
}
func (e errno) GetNowTime() int64 {
	return e.NowTime
}
func (e errno) GetRequestId() string {
	return e.RequestID
}

func (e errno) Error() string {
	return e.Msg
}

func (e errno) GetCode() int {
	return e.Code
}

func (e errno) GetHttpStatusCode() int {
	e.resetHttpStatusCode()
	return e.HttpStatusCode
}

func (e errno) GetErrMsg() string {
	return e.Msg
}

func (e errno) GetData() interface{} {
	return e.Data
}

// 仅支持string/[]byte类型的data
func (e errno) GetRawData() []byte {
	if s, ok := e.Data.(string); ok {
		return []byte(s)
	}

	if s, ok := e.Data.([]byte); ok {
		return s
	}

	return []byte{}
}

// DEPRECATED 请使用 FormatErrMsg
func (e errno) Format(v ...interface{}) Error {
	e.Msg = fmt.Sprintf(e.Msg, v...)

	return e
}

func (e errno) FormatErrMsg(v ...interface{}) Error {
	e.Msg = fmt.Sprintf(e.Msg, v...)

	return e
}

func (e errno) WithData(data interface{}) Error {
	e.Data = data
	return e
}

func (e errno) WithID(rid string) Error {
	e.RequestID = rid

	return e
}

func (e errno) reset(ctx *gin.Context) Error {
	e.NowTime = time.Now().UnixMilli()
	if e.GetCode() == OK.GetCode() {
		e.resetHttpStatusCode()
	}
	return e.requestId(ctx)
}

func (e errno) requestId(ctx context.Context) Error {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if span := trace.SpanContextFromContext(ginCtx.Request.Context()); span.HasTraceID() {
			e.RequestID = span.TraceID().String()
		}
	}
	return e.WithHttpStatusCode(e.Code)
	//return e
}

func (e *errno) resetHttpStatusCode() {
	if e.GetCode() == OK.GetCode() {
		e.HttpStatusCode = http.StatusOK
	} else if e.GetCode() == ForbiddenTimeOut.GetCode() {
		e.HttpStatusCode = http.StatusUnauthorized
	} else if e.GetCode() == Forbidden.GetCode() {
		e.HttpStatusCode = http.StatusForbidden
	} else if e.GetCode() == SysErr.GetCode() || e.GetCode() == NetworkErr.GetCode() {
		e.HttpStatusCode = http.StatusInternalServerError
	} else if e.GetCode() == RecordNotFound.GetCode() {
		e.HttpStatusCode = http.StatusNotFound
	} else {
		e.HttpStatusCode = http.StatusBadRequest
	}
}

func (e errno) WithGrpcError(err error) Error {
	rpcError := errors.FromError(err)
	e.resetHttpStatusCode()
	e.Msg = rpcError.GetMessage()
	e.ErrorCode = rpcError.GetReason()
	e.Reason = rpcError.GetReason()
	e.Code = int(rpcError.GetCode())
	return e
}

// 如果reason是Error，则会直接使用reason返回
func (e errno) WithReason(reason interface{}) Error {
	if v, ok := reason.(Error); ok {
		e.resetHttpStatusCode()
		return v.WithID(e.RequestID)
	}

	if v, ok := reason.(error); ok {
		e.resetHttpStatusCode()
		e.Msg = v.Error()
		e.Reason = v.Error()
		return e
	}
	e.resetHttpStatusCode()
	if e.Code == ParamValidationErr.GetCode() {
		e.Msg = fmt.Sprintf("%v", reason)
	}
	e.Reason = reason
	return e
}

// 请使用net.http status code
func (e errno) WithHttpStatusCode(httpStatusCode int) Error {
	e.HttpStatusCode = httpStatusCode
	if e.GetCode() == 200 {
		e.HttpStatusCode = http.StatusOK
	} else if e.GetCode() >= 111 && e.GetCode() <= 112 {
		e.HttpStatusCode = http.StatusUnauthorized
	} else {
		e.HttpStatusCode = http.StatusBadRequest
	}
	return e
}

func (e errno) ToString() string {
	return string(e.ToBytes())
}

func (e errno) ToDataString() string {
	return string(e.ToDataBytes())
}

func (e errno) ToDataBytes() []byte {
	data := e.Data
	if data == nil {
		data = make(map[string]interface{})
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return raw
}

func (e errno) ToBytes() []byte {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}

	raw, err := json.Marshal(e)
	if err != nil {
		return nil
	}

	return raw
}
