package field

import (
	"github.com/tinylib/msgp/msgp"
)

func (b *Bool) DecodeMsg(r *msgp.Reader) error {
	value, err := r.ReadBool()
	if err != nil {
		return err
	}
	return b.Scan(value)
}

func (b *Bool) MarshalMsg(i []byte) ([]byte, error) {
	o := msgp.Require(i, b.Msgsize())
	o = msgp.AppendBool(o, b.Bool)
	return o, nil
}

func (b *Bool) UnmarshalMsg(i []byte) ([]byte, error) {
	value, o, e := msgp.ReadBoolBytes(i)
	if e != nil {
		return o, e
	}
	return o, b.Scan(value)
}

func (b *Bool) Msgsize() int {
	return msgp.BoolSize
}

func (b *Bool) EncodeMsg(w *msgp.Writer) error {
	return w.WriteBool(b.Bool)
}

func (nb *NullBool) EncodeMsg(w *msgp.Writer) error {
	if nb.Valid {
		return w.WriteBool(nb.Bool)
	}
	return w.WriteNil()
}

func (nb *NullBool) DecodeMsg(r *msgp.Reader) error {
	if r.IsNil() {
		err := r.ReadNil()
		if err != nil {
			return err
		}
		return nb.Scan(nil)
	}
	value, err := r.ReadBool()
	if err != nil {
		return err
	}
	return nb.Scan(value)

}

func (nb *NullBool) MarshalMsg(i []byte) ([]byte, error) {
	if !nb.Valid {
		o := msgp.Require(i, msgp.NilSize)
		o = msgp.AppendNil(o)
		return o, nil
	}
	o := msgp.Require(i, nb.Msgsize())
	o = msgp.AppendBool(o, nb.Bool)
	return o, nil

}

func (nb *NullBool) UnmarshalMsg(i []byte) ([]byte, error) {
	t := msgp.NextType(i)
	if t == msgp.NilType {
		o, err := msgp.ReadNilBytes(i)
		if err != nil {
			return o, err
		}
		return o, nb.Scan(nil)
	}
	value, o, err := msgp.ReadBoolBytes(i)
	if err != nil {
		return o, err
	}
	return o, nb.Scan(value)
}

func (nb *NullBool) Msgsize() int {
	if !nb.Valid {
		return msgp.NilSize
	}
	return msgp.BoolSize
}
