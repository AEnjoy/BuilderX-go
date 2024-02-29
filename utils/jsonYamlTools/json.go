package jsonYamlTools

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

// GetFieldValue 获取嵌套结构体中的值(参数1：字段路径，参数2：解析后的数据)
func GetFieldValue(path []string, data map[string]interface{}) (interface{}, bool) {
	if v, ok := data[path[0]]; ok == true {
		if len(path) == 1 {
			return v, true
		} else {
			value, ok := GetFieldValue(path[1:], v.(map[string]interface{}))
			return value, ok
		}
	} else {
		return nil, false
	}
}

// GetFieldFromJson 获取指定字段的值(参数1：字段路径，参数2：原始json数据)
func GetFieldFromJson(path []string, value []byte) (string, error) {
	var temp *ast.Node
	for i, v := range path {
		if i == 0 {
			_, err := sonic.Get(value, v)
			if err != nil {
				return "", err
			}
		} else {
			temp = temp.Get(v)
			if temp != nil {
				return "", errors.New("path not found")
			}
		}
	}
	switch temp.Type() {
	case ast.V_NULL:
		return "", errors.New("path not found")
	case ast.V_STRING:
		return temp.String()
	case ast.V_NUMBER:
		t, err := temp.Number()
		return t.String(), err
	}
	return "", errors.New("path not found")
}
