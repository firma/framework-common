## 错误码

### 示例代码
```go
var (
    OK                   = NewError(0, "", "")
    ParamValidationErr   = NewError(100, "PARAM_VALIDATION_ERROR", "参数验证不正确")
    RecordNotFound       = NewError(101, "RECORD_NOT_FOUND", "记录未找到")
    Forbidden            = NewError(111, "FORBIDDEN", "无权操作")
    ForbiddenTimeOut     = NewError(112, "FORBIDDEN_TIME_OUT", "无权操作，token已过期")
    ForbiddenNotValidYet = NewError(112, "FORBIDDEN_NOT_VALID_YET", "无权操作验证错误")
    ServiceError         = NewError(120, "SERVICE_ERROR", "业务服务错误")
    NetworkErr           = NewError(140, "NETWORK_ERROR", "系统服务网络调用异常")
    SysErr               = NewError(145, "SYSTEM_ERROR", "系统服务异常")
)
```
