package jq

/*
#cgo LDFLAGS: -l jq
#include <jq.h>
#include <stdlib.h>

static void on_jq_err(void *cx, jv err) {
	if (*(const char**) cx != NULL) {
		return;
	}

	*(const char**) cx = jv_string_value(err);
}

static void set_err_cb(jq_state *jq, const char **msg) {
	if (msg == NULL) {
		jq_set_error_cb(jq, NULL, NULL);
		return;
	}

	jq_set_error_cb(jq, on_jq_err, msg);
}
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Run will run the program over the given value input. The JQ parse
// error is returned if the given expression fails compilation.
func Run(program string, value interface{}) ([]interface{}, error) {
	jq := C.jq_init()
	defer C.jq_teardown(&jq)

	if err := compile(jq, program); err != nil {
		return nil, err
	}

	C.jq_start(jq, marshal(value), 0)
	return readValues(jq), nil
}

func compile(jq *C.jq_state, program string) error {
	cs := C.CString(program)
	defer C.free(unsafe.Pointer(cs))

	var errMessage *C.char
	C.set_err_cb(jq, &errMessage)
	defer C.set_err_cb(jq, nil)

	if C.jq_compile(jq, cs) == 0 {
		return fmt.Errorf(C.GoString(errMessage))
	}

	return nil
}

func readValues(jq *C.jq_state) []interface{} {
	values := []interface{}{}

	for {
		value, ok := readValue(jq)
		if !ok {
			break
		}

		values = append(values, value)
	}

	return values
}

func readValue(jq *C.jq_state) (interface{}, bool) {
	value := C.jq_next(jq)
	defer C.jv_free(value)

	if C.jv_is_valid(value) == 0 {
		return nil, false
	}

	return unmarshal(value), true
}

func marshal(v interface{}) C.jv {
	if v == nil {
		return C.jv_null()
	}

	value := reflect.Indirect(reflect.ValueOf(v))

	switch value.Type().Kind() {
	case reflect.Bool:
		return marshalBool(value)
	case reflect.String:
		return C.jv_string(C.CString(value.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return C.jv_number(C.double(value.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return C.jv_number(C.double(value.Uint()))
	case reflect.Float32, reflect.Float64:
		return C.jv_number(C.double(value.Float()))
	case reflect.Slice, reflect.Array:
		return marshalArray(value)
	case reflect.Map:
		return marshalObject(value)
	default:
		errMessage := fmt.Sprintf("unknown type %v", value.Interface())
		return C.jv_invalid_with_msg(C.jv_string(C.CString(errMessage)))
	}
}

func marshalBool(value reflect.Value) C.jv {
	if value.Bool() {
		return C.jv_true()
	}

	return C.jv_false()
}

func marshalArray(value reflect.Value) C.jv {
	jvArray := C.jv_array_sized(C.int(value.Len()))

	for i := 0; i < value.Len(); i++ {
		jvArray = C.jv_array_set(
			C.jv_copy(jvArray),
			C.int(i),
			marshal(value.Index(i).Interface()),
		)
	}

	return jvArray
}

func marshalObject(value reflect.Value) C.jv {
	jvObject := C.jv_object()

	for _, k := range value.MapKeys() {
		jvObject = C.jv_object_set(
			jvObject,
			marshal(k.Interface()),
			marshal(value.MapIndex(k).Interface()),
		)
	}

	return jvObject
}

func unmarshal(value C.jv) interface{} {
	switch C.jv_get_kind(value) {
	case C.JV_KIND_NULL:
		return nil
	case C.JV_KIND_TRUE:
		return true
	case C.JV_KIND_FALSE:
		return false
	case C.JV_KIND_STRING:
		return C.GoString(C.jv_string_value(value))
	case C.JV_KIND_NUMBER:
		return unmarshalNumber(value)
	case C.JV_KIND_ARRAY:
		return unmarshalArray(value)
	case C.JV_KIND_OBJECT:
		return unmarshalObject(value)
	case C.JV_KIND_INVALID:
		return fmt.Errorf("invalid value")
	default:
		return fmt.Errorf("unknown value")
	}
}

func unmarshalNumber(value C.jv) interface{} {
	if C.jv_is_integer(value) == 0 {
		return float64(C.jv_number_value(value))
	}

	return int(C.jv_number_value(value))
}

func unmarshalArray(value C.jv) interface{} {
	s := []interface{}{}
	n := int(C.jv_array_length(C.jv_copy(value)))

	for i := 0; i < n; i++ {
		s = append(s, unmarshal(C.jv_array_get(C.jv_copy(value), C.int(i))))
	}

	return s
}

func unmarshalObject(value C.jv) interface{} {
	m := map[string]interface{}{}
	for iter := C.jv_object_iter(value); C.jv_object_iter_valid(value, iter) != 0; iter = C.jv_object_iter_next(value, iter) {
		rk := C.jv_object_iter_key(value, iter)
		sk := C.jv_string_value(rk)
		rv := C.jv_object_iter_value(value, iter)

		m[C.GoString(sk)] = unmarshal(rv)
	}

	return m
}

// Exposed for tests
func free(value C.jv) {
	C.jv_free(value)
}
