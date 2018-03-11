package field

import (
	"github.com/tinylib/msgp/msgp"
)

func (b *Float64) DecodeMsg(r *msgp.Reader) error {
	value, err := r.ReadFloat64()
	if err != nil {
		return err
	}
	return b.Scan(value)
}

func (b *Float64) MarshalMsg(i []byte) ([]byte, error) {
	o := msgp.Require(i, b.Msgsize())
	o = msgp.AppendFloat64(o, b.Float64)
	return o, nil
}

func (b *Float64) UnmarshalMsg(i []byte) ([]byte, error) {
	value, o, e := msgp.ReadFloat64Bytes(i)
	if e != nil {
		return o, e
	}
	return o, b.Scan(value)
}

func (b *Float64) Msgsize() int {
	return msgp.Float64Size
}

func (b *Float64) EncodeMsg(w *msgp.Writer) error {
	return w.WriteFloat64(b.Float64)
}

func (nb *NullFloat64) EncodeMsg(w *msgp.Writer) error {
	if nb.Valid {
		return w.WriteFloat64(nb.Float64)
	}
	return w.WriteNil()
}

func (nb *NullFloat64) DecodeMsg(r *msgp.Reader) error {
	if r.IsNil() {
		err := r.ReadNil()
		if err != nil {
			return err
		}
		return nb.Scan(nil)
	}
	value, err := r.ReadFloat64()
	if err != nil {
		return err
	}
	return nb.Scan(value)

}

func (nb *NullFloat64) MarshalMsg(i []byte) ([]byte, error) {
	if !nb.Valid {
		o := msgp.Require(i, msgp.NilSize)
		o = msgp.AppendNil(o)
		return o, nil
	}
	o := msgp.Require(i, nb.Msgsize())
	o = msgp.AppendFloat64(o, nb.Float64)
	return o, nil

}

func (nb *NullFloat64) UnmarshalMsg(i []byte) ([]byte, error) {
	t := msgp.NextType(i)
	if t == msgp.NilType {
		o, err := msgp.ReadNilBytes(i)
		if err != nil {
			return o, err
		}
		return o, nb.Scan(nil)
	}
	value, o, err := msgp.ReadFloat64Bytes(i)
	if err != nil {
		return o, err
	}
	return o, nb.Scan(value)
}

func (nb *NullFloat64) Msgsize() int {
	if !nb.Valid {
		return msgp.NilSize
	}
	return msgp.Float64Size
}
