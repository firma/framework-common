## 验证器(仅作用于gin)

### 示例代码
```
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		return errno.ParamValidationErr.WithReason(validate.TransErr(err))
	}
```