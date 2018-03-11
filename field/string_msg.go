package field

import (
	"github.com/tinylib/msgp/msgp"
)

func (b *String) DecodeMsg(r *msgp.Reader) error {
	value, err := r.ReadString()
	if err != nil {
		return err
	}
	return b.Scan(value)
}

func (b *String) MarshalMsg(i []byte) ([]byte, error) {
	o := msgp.Require(i, b.Msgsize())
	o = msgp.AppendString(o, b.String)
	return o, nil
}

func (b *String) UnmarshalMsg(i []byte) ([]byte, error) {
	value, o, e := msgp.ReadStringBytes(i)
	if e != nil {
		return o, e
	}
	return o, b.Scan(value)
}

func (b *String) Msgsize() int {
	return msgp.StringPrefixSize + len(b.String)
}

func (b *String) EncodeMsg(w *msgp.Writer) error {
	return w.WriteString(b.String)
}

func (nb *NullString) EncodeMsg(w *msgp.Writer) error {
	if nb.Valid {
		return w.WriteString(nb.String)
	}
	return w.WriteNil()
}

func (nb *NullString) DecodeMsg(r *msgp.Reader) error {
	if r.IsNil() {
		err := r.ReadNil()
		if err != nil {
			return err
		}
		return nb.Scan(nil)
	}
	value, err := r.ReadString()
	if err != nil {
		return err
	}
	return nb.Scan(value)

}

func (nb *NullString) MarshalMsg(i []byte) ([]byte, error) {
	if !nb.Valid {
		o := msgp.Require(i, msgp.NilSize)
		o = msgp.AppendNil(o)
		return o, nil
	}
	o := msgp.Require(i, nb.Msgsize())
	o = msgp.AppendString(o, nb.String)
	return o, nil

}

func (nb *NullString) UnmarshalMsg(i []byte) ([]byte, error) {
	t := msgp.NextType(i)
	if t == msgp.NilType {
		o, err := msgp.ReadNilBytes(i)
		if err != nil {
			return o, err
		}
		return o, nb.Scan(nil)
	}
	value, o, err := msgp.ReadStringBytes(i)
	if err != nil {
		return o, err
	}
	return o, nb.Scan(value)
}

func (nb *NullString) Msgsize() int {
	if !nb.Valid {
		return msgp.NilSize
	}
	return msgp.StringPrefixSize + len(nb.String)
}
