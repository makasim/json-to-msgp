package jsontomsgp

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"
	"github.com/valyala/fastjson"
)

var pp = &fastjson.ParserPool{}

func CopyBytes(src []byte, w *msgp.Writer) error {
	p := pp.Get()
	defer pp.Put(p)

	v, err := p.ParseBytes(src)
	if err != nil {
		return err
	}

	if err := convertV(v, w); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func convertV(v *fastjson.Value, w *msgp.Writer) error {
	defer w.Flush()

	switch v.Type() {
	case fastjson.TypeObject:
		objV, err := v.Object()
		if err != nil {
			return err
		}

		objV.Len()

		err = w.WriteMapHeader(uint32(objV.Len()))
		if err != nil {
			return err
		}

		var visitErr error
		objV.Visit(func(key []byte, v *fastjson.Value) {
			if visitErr != nil {
				return
			}

			if err := w.WriteStringFromBytes(key); err != nil {
				visitErr = err
			}
			if err := convertV(v, w); err != nil {
				visitErr = err
			}
		})

		if visitErr != nil {
			return visitErr
		}
		return nil
	case fastjson.TypeArray:
		items, err := v.Array()
		if err != nil {
			return err
		}

		sz := uint32(len(items))
		err = w.WriteArrayHeader(sz)
		if err != nil {
			return err
		}
		for i := range items {
			if err := convertV(items[i], w); err != nil {
				return err
			}
		}
		return nil
	case fastjson.TypeString:
		if err := w.WriteStringFromBytes(v.GetStringBytes()); err != nil {
			return err
		}
		return nil
	case fastjson.TypeNumber:
		if err := w.WriteFloat64(v.GetFloat64()); err != nil {
			return err
		}
		return nil
	case fastjson.TypeTrue:
		if err := w.WriteBool(true); err != nil {
			return err
		}
		return nil
	case fastjson.TypeFalse:
		if err := w.WriteBool(false); err != nil {
			return err
		}
		return nil
	case fastjson.TypeNull:
		if err := w.WriteNil(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown type: %s", v.Type())
	}
}
