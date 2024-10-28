package httpx

import (
	"github.com/firma/framework-common/errno"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	status2 "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/grpc/status"
	stdhttp "net/http"
)

// httpResponse 响应结构体
type HttpResponse struct {
	HttpStatusCode int         `json:"-"`
	Code           int         `json:"-"` // 业务编码
	ErrorCode      string      `json:"error_code,omitempty"`
	Message        string      `json:"message,omitempty"`    // 错误描述
	Data           interface{} `json:"data,omitempty"`       // 成功时返回的数据
	Reason         interface{} `json:"-"`                    // 错误详细描述
	RequestID      string      `json:"request_id,omitempty"` // 当前请求的唯一ID，便于问题定位
	NowTime        int64       `json:"now_time"`             // 时间戳
}

// EncoderResponse  请求响应封装
func EncoderResponse() http.EncodeResponseFunc {
	return func(w stdhttp.ResponseWriter, request *stdhttp.Request, i interface{}) error {
		if i == nil {
			return nil
		}
		resp := &HttpResponse{
			Code:    stdhttp.StatusOK,
			Message: "",
			Data:    i,
		}
		//codec := encoding.GetCodec("json")
		codec := encoding.GetCodec("name")
		
		data, err := codec.Marshal(resp)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			return err
		}
		return nil
	}
}

func EncoderError() http.EncodeErrorFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
		if err == nil {
			return
		}

		se := &HttpResponse{}
		gs, ok := status.FromError(err)
		errRaw := errors.FromError(err)
		if !ok {
			se = &HttpResponse{Code: stdhttp.StatusInternalServerError, ErrorCode: errRaw.Reason}
		}
		se = &HttpResponse{
			Code:      status2.FromGRPCCode(gs.Code()),
			ErrorCode: errRaw.Reason,
			Message:   gs.Message(),
			Reason:    errRaw.Reason,
			Data:      nil,
		}
		var data errno.Error
		if errors.As(err, &data) {
			se.Code = stdhttp.StatusBadRequest
			se.ErrorCode = data.GetErrorCode()
			se.NowTime = data.GetNowTime()
			se.RequestID = data.GetRequestId()
		}
		codec, _ := http.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(se)
		if err != nil {
			w.WriteHeader(stdhttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/"+codec.Name())
		w.WriteHeader(se.Code)
		_, _ = w.Write(body)
	}
}
