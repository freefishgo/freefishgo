package mvc

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// 字符转类型
func MapStringToStruct(i interface{}, data map[string]interface{}) interface{} {
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		if !v.Elem().CanSet() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		v1 := v.Field(i)
		if name == "" {
			name = f.Name
			_, ok := data[name]
			if !ok {
				if v1.Kind() == reflect.Struct || (v1.Kind() == reflect.Ptr && v1.Elem().Kind() == reflect.Struct) {
					MapStringToStructInReflect(v1, data)
				}
				continue
			}
		}
		val, ok := data[name]
		if !ok {
			continue
		}
		switch v1.Kind() {
		case reflect.Slice:
			if val1, ok := val.([]string); ok {
				sl := reflect.MakeSlice(v1.Type(), len(val1), len(val1))
				b := false
				for i := 0; i < len(val1); i++ {
					if doBasic(sl.Index(i), f.Type.Elem(), val1[i]) {
						b = true
					}
				}
				if b {
					v1.Set(sl)
				}
			} else if val2, ok := val.(string); ok {
				vTemp := reflect.New(v1.Type()).Interface()
				err := json.Unmarshal([]byte(val2), vTemp)
				if err == nil {
					v1.Set(reflect.ValueOf(vTemp).Elem())
				}
			}
			continue
		case reflect.Ptr:
			if !v1.Elem().CanSet() {
				v1.Set(reflect.New(f.Type.Elem()))
			}
			doBasic(v1.Elem(), f.Type.Elem(), val)
			continue
		default:
			doBasic(v1, f.Type, val)
			continue
		}
	}
	return i
}

func MapStringToStructInReflect(v reflect.Value, data map[string]interface{}) interface{} {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		if !v.Elem().CanSet() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		v1 := v.Field(i)
		if name == "" {
			name = f.Name
			_, ok := data[name]
			if !ok {
				if v1.Kind() == reflect.Struct || (v1.Kind() == reflect.Ptr && v1.Elem().Kind() == reflect.Struct) {
					MapStringToStructInReflect(v1, data)
				}
				continue
			}
		}
		val, ok := data[name]
		if !ok {
			continue
		}
		switch v1.Kind() {
		case reflect.Slice:
			if val1, ok := val.([]string); ok {
				sl := reflect.MakeSlice(v1.Type(), len(val1), len(val1))
				b := false
				for i := 0; i < len(val1); i++ {
					if doBasic(sl.Index(i), f.Type.Elem(), val1[i]) {
						b = true
					}
				}
				if b {
					v1.Set(sl)
				}
			} else if val2, ok := val.(string); ok {
				vTemp := reflect.New(v1.Type()).Interface()
				err := json.Unmarshal([]byte(val2), vTemp)
				if err == nil {
					v1.Set(reflect.ValueOf(vTemp).Elem())
				}
			}
			continue
		case reflect.Ptr:
			if !v1.Elem().CanSet() {
				v1.Set(reflect.New(f.Type.Elem()))
			}
			doBasic(v1.Elem(), f.Type.Elem(), val)
			continue
		default:
			doBasic(v1, f.Type, val)
			continue
		}
	}
	return v
}

func doBasic(v1 reflect.Value, t1 reflect.Type, val interface{}) bool {
	switch t1.Kind() {
	case reflect.String:
		if val, ok := val.(string); ok {
			if v1.CanSet() {
				v1.SetString(val)
			}
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, ok := val.(string); ok {
			if val, err := strconv.ParseInt(val, 10, v1.Type().Bits()); err == nil {
				if v1.CanSet() {
					v1.SetInt(val)
				}
			}
		}
		break
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, ok := val.(string); ok {
			if val, err := strconv.ParseUint(val, 10, v1.Type().Bits()); err == nil {
				if v1.CanSet() {
					v1.SetUint(val)
				}
			}
		}
		break
	case reflect.Bool:
		if val, ok := val.(string); ok {
			if val, err := strconv.ParseBool(val); err == nil {
				if v1.CanSet() {
					v1.SetBool(val)
				}
			}
		}
		break
	case reflect.Float32, reflect.Float64:
		if val, ok := val.(string); ok {
			if val, err := strconv.ParseFloat(val, v1.Type().Bits()); err == nil {
				if v1.CanSet() {
					v1.SetFloat(val)
				}
			}
		}
		break
	case reflect.Ptr:
		t1 = t1.Elem()
		if !v1.Elem().CanSet() {
			v1.Set(reflect.New(t1))
		}
		v1 = v1.Elem()
		doBasic(v1, t1, val)
		break
	case reflect.Struct:
		if val, ok := val.(string); ok {
			vTemp := reflect.New(t1).Interface()
			err := json.Unmarshal([]byte(val), vTemp)
			if err == nil {
				v1.Set(reflect.ValueOf(vTemp).Elem())
			}
		}
	case reflect.Interface:
		if v1.CanSet() {
			v1.Set(reflect.ValueOf(val))
		}
	default:
		return false
	}
	return true
}
