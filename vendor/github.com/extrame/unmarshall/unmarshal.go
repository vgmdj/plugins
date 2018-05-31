package unmarshall

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

//Unmarshaller for specified struct, ValueGetter for how to get the value, it get the tag as input and return the value
type Unmarshaller struct {
	ValueGetter  func(string) []string
	TagConcatter func(string, string) string
	// FileGetter            func(string) (multipart.File, *multipart.FileHeader, error)
	FillForSpecifiledType map[string]func(string) (reflect.Value, error)
	AutoFill              bool
}

func (u *Unmarshaller) Unmarshall(v interface{}) error {
	// check v is valid
	rv := reflect.ValueOf(v).Elem()
	// dereference pointer
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		// for each struct field on v
		u.unmarshalStructInForm("", rv, 0, false, make(map[string]bool))
	} else {
		return fmt.Errorf("v must point to a struct type")
	}
	return nil
}

func (u *Unmarshaller) unmarshalStructInForm(context string,
	rvalue reflect.Value,
	offset int,
	inarray bool,
	parents map[string]bool) (err error) {

	if rvalue.Type().Kind() == reflect.Ptr {

		rvalue = rvalue.Elem()
	}
	rtype := rvalue.Type()

	parents[rtype.PkgPath()+"/"+rtype.Name()] = true

	success := false

	for i := 0; i < rtype.NumField() && err == nil; i++ {
		id, form_values, tag := u.getFormField(context, rtype.Field(i), offset, inarray)
		increaseOffset := !(context != "" && inarray)
		var used_offset = 0
		if increaseOffset {
			used_offset = offset
		}
		if rvalue.Field(i).CanSet() {
			switch rtype.Field(i).Type.Kind() {
			case reflect.Ptr: //TODO if the ptr point to a basic data, it will crash
				val := rvalue.Field(i)
				typ := rtype.Field(i).Type.Elem()
				if val.IsNil() {
					val.Set(reflect.New(typ))
				}
				if err = u.fill_struct(typ, val.Elem(), id, form_values, tag, used_offset, parents); err != nil {
					return err
				} else {
					break
				}
			case reflect.Struct:
				if err = u.fill_struct(rtype.Field(i).Type, rvalue.Field(i), id, form_values, tag, used_offset, parents); err != nil {
					glog.Errorln(err)
					return err
				} else {
					break
				}
			case reflect.Interface:
				//ask the parent to tell me how to unmarshal it
				mtd := rvalue.MethodByName("UnmarshallForm")
				if mtd.IsValid() {
					values := mtd.Call([]reflect.Value{reflect.ValueOf(rtype.Field(i).Name)})
					if len(values) == 2 && values[1].Interface() == nil {
						res := values[0].Interface()
						resValue := reflect.ValueOf(res)
						resType := reflect.TypeOf(res)
						if err = u.fill_struct(resType, resValue, id, form_values, tag, used_offset, parents); err != nil {
							rvalue.Field(i).Set(resValue)
							return err
						} else {
							break
						}
					}
				} else {
					glog.Infoln(fmt.Errorf("try to use UnmarshallForm to unmarshall interface type(%T) fail", rvalue.Interface()))
				}
			case reflect.Slice:
				fType := rtype.Field(i).Type
				subRType := rtype.Field(i).Type.Elem()
				if fType.PkgPath() == "net" && fType.Name() == "IP" && len(form_values) > 0 && used_offset < len(form_values) {
					rvalue.Field(i).Set(reflect.ValueOf(net.ParseIP(form_values[used_offset])))
					continue
				}
				switch subRType.Kind() {
				case reflect.Struct:
					if _, ok := parents[subRType.PkgPath()+"/"+subRType.Name()]; !ok {
						rvalueTemp := reflect.MakeSlice(rtype.Field(i).Type, 0, 0)
						subRValue := reflect.New(subRType)
						offset := 0
						for {
							err = u.unmarshalStructInForm(id, subRValue, offset, true, parents)
							if err != nil {
								break
							}
							offset++
							rvalueTemp = reflect.Append(rvalueTemp, subRValue.Elem())
						}
						rvalue.Field(i).Set(rvalueTemp)
					} else {
						err = fmt.Errorf("Too deep of type reuse %v", parents)
					}
				default:
					len_fv := len(form_values)
					rvnew := reflect.MakeSlice(rtype.Field(i).Type, len_fv, len_fv)
					for j := 0; j < len_fv; j++ {
						u.unmarshalField(context, rvnew.Index(j), form_values[j], tag)
					}
					rvalue.Field(i).Set(rvnew)
				}
			default:
				if len(form_values) > 0 && used_offset < len(form_values) {
					u.unmarshalField(context, rvalue.Field(i), form_values[used_offset], tag)
					success = true
				} else if len(tag) > 0 {
					u.unmarshalField(context, rvalue.Field(i), tag[0], tag)
				}
			}
		} else {
			glog.Errorf("cannot set value of (%s,%s) in fill", rtype.Field(i).Name, rtype.Field(i).Type.Name())
		}
	}
	if !success && err == nil {
		err = errors.New("no more element")
	}
	return
}

func (u *Unmarshaller) getFormField(prefix string, t reflect.StructField, offset int, inarray bool) (string, []string, []string) {

	tag, tags := u.getTag(prefix, t, offset, inarray)

	values := u.ValueGetter(tag)

	return tag, values, tags[1:]
}

func (u *Unmarshaller) getTag(prefix string,
	t reflect.StructField, offset int, inarray bool) (string, []string) {
	tags := []string{""}
	tag := t.Tag.Get("goblet")
	if tag != "" {
		tags = strings.Split(tag, ",")
		tag = tags[0]
	}
	if tag == "" {
		tag = t.Name
	}

	// values := []string{}

	// if form != nil {
	// 	values = (*form)[tag]
	// }

	if prefix != "" {
		if inarray {
			tag = fmt.Sprintf(prefix+"[%d]"+"["+tag+"]", offset)
		} else {
			tag = u.TagConcatter(prefix, tag)
		}
	}
	return tag, tags
}

func (u *Unmarshaller) unmarshalField(contex string, v reflect.Value, form_value string, tags []string) error {
	// string -> type conversion
	switch v.Kind() {
	case reflect.Int64:
		if i, err := strconv.ParseInt(form_value, 10, 64); err == nil {
			v.SetInt(i)
		}
	case reflect.Uint64:
		if i, err := strconv.ParseUint(form_value, 10, 64); err == nil {
			v.SetUint(i)
		}
	case reflect.Int, reflect.Int32:
		if i, err := strconv.ParseInt(form_value, 10, 32); err == nil {
			v.SetInt(i)
		}
	case reflect.Uint32:
		if i, err := strconv.ParseUint(form_value, 10, 32); err == nil {
			v.SetUint(i)
		}
	case reflect.Int16:
		if i, err := strconv.ParseInt(form_value, 10, 16); err == nil {
			v.SetInt(i)
		}
	case reflect.Uint16:
		if i, err := strconv.ParseUint(form_value, 10, 16); err == nil {
			v.SetUint(i)
		}
	case reflect.Int8:
		if i, err := strconv.ParseInt(form_value, 10, 8); err == nil {
			v.SetInt(i)
		}
	case reflect.Uint8:
		if i, err := strconv.ParseUint(form_value, 10, 8); err == nil {
			v.SetUint(i)
		}
	case reflect.String:
		// copy string
		if len(tags) > 0 && tags[len(tags)-1] == "md5" {
			h := md5.New()
			h.Write([]byte(form_value))
			v.SetString(hex.EncodeToString(h.Sum(nil)))
		} else {
			v.SetString(form_value)
		}
	case reflect.Float64:
		if f, err := strconv.ParseFloat(form_value, 64); err == nil {
			v.SetFloat(f)
		}
	case reflect.Float32:
		if f, err := strconv.ParseFloat(form_value, 32); err == nil {
			v.SetFloat(f)
		}
	case reflect.Bool:
		// the following strings convert to true
		// 1,true,on,yes
		fv := form_value
		if fv == "1" || fv == "true" || fv == "on" || fv == "yes" {
			v.SetBool(true)
		}
	default:
		fmt.Println("unknown type", v.Kind())
	}
	return nil
}

func (u *Unmarshaller) fill_struct(typ reflect.Type,
	val reflect.Value, id string, form_values []string, tag []string, used_offset int, parents map[string]bool) error {
	if typ.PkgPath() == "time" && typ.Name() == "Time" {
		var fillby string
		var fillby_valid = regexp.MustCompile(`^\s*fillby\((.*)\)\s*$`)
		for _, v := range tag {
			matched := fillby_valid.FindStringSubmatch(v)
			if len(matched) == 2 {
				fillby = matched[1]
			}
		}
		fillby = strings.TrimSpace(fillby)
		var value string
		if len(form_values) > used_offset {
			value = form_values[used_offset]
		}

		switch fillby {
		case "now":
			val.Set(reflect.ValueOf(time.Now()))
		case "timestamp":
			if unix, err := strconv.ParseInt(value, 10, 64); err == nil {
				val.Set(reflect.ValueOf(time.Unix(unix, 0)))
			} else {
				return err
			}
		default:
			if fillby == "" {
				fillby = time.RFC3339
			}
			if value != "" {
				time, err := time.Parse(fillby, value)
				if err == nil {
					val.Set(reflect.ValueOf(time))
				} else {
					log.Println(err)
					return err
				}
			}
		}
	} else {
		for k, fn := range u.FillForSpecifiledType {
			if typ.PkgPath()+"."+typ.Name() == k {
				if v, err := fn(id); err == nil {
					val.Set(v)
					return nil
				} else {
					return err
				}
			}
		}
		u.unmarshalStructInForm(id, val, 0, false, parents)
	}
	return nil
}
