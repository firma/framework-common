package domain

import "encoding/json"

func ParamsConvert[T any, M any](t *T, m M) *M {
	if t == nil {
		return nil
	}
	if dataByte, err := json.Marshal(t); err != nil {
		return nil
	} else {
		err := json.Unmarshal(dataByte, &m)
		if err != nil {
			return nil
		}
		return &m
	}
}
