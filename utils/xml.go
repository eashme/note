/*
*	@Time : 2020-7-10 03:05 下午
*   @Author : jake
*	接入微信服务号消息响应 进行 xml 编码
*	没找到 合适的xml解压方式
*	根据微信服务号文档要求的xml格式用反射实现了一个
*	因为完全使用反射,也没做缓存,性能应该挺差的,暂且堪用
*/
package utils

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func Marshal(v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v)

	buf := bytes.NewBuffer([]byte{})

	_, err := buf.Write([]byte("<xml>")) // 写头
	if err != nil {
		return nil, err
	}

	err = encodeObj(buf, val)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte("</xml>")) // 写尾
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func encodeObj(w io.Writer, v reflect.Value) error {

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.IsZero() {
		return nil
	}

	typ := v.Type()
	n := typ.NumField()
	for i := 0; i < n; i++ {
		// 获取字段
		field := typ.Field(i)
		// 获取类型
		value := v.Field(i)
		// 对单个字段进行编码
		err := encodeField(w, field, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// 仅支持 字符串 , 数字 , 浮点 , 指针, 结构体 , Slice 类型 其他类型不支持 [指针,结构体 类型中包含]
func encodeField(w io.Writer, f reflect.StructField, v reflect.Value) error {

	// 0值直接返回
	if v.IsZero() {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 头标签
	tagName, ok := f.Tag.Lookup("xml")
	if !ok {
		tagName = f.Name
	}

	_, err := w.Write([]byte("<" + tagName + ">"))
	if err != nil {
		return err
	}

	k := v.Kind()
	switch k {
	case reflect.String: // string类型需要使用 ![CDATA[xxx]] 进行包裹
		_, err = w.Write([]byte("![CDATA[" + v.String() + "]]"))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err = w.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err = w.Write([]byte(strconv.FormatUint(v.Uint(), 10)))
	case reflect.Float64, reflect.Float32:
		_, err = w.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
	case reflect.Struct: // 嵌套
		err = encodeObj(w, v)
	case reflect.Slice:
		n := v.Len()
		for i := 0; i < n; i++ {
			_, err = w.Write([]byte("<item>"))
			if err != nil {
				return err
			}
			err = encodeObj(w, v.Index(i))
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("</item>"))
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("not support %s", k.String())
	}
	if err != nil {
		return err
	}

	// 结尾标签
	_, err = w.Write([]byte("</" + tagName + ">"))
	if err != nil {
		return err
	}

	return nil
}
