package field

import (
	"github.com/tinylib/msgp/msgp"
)

func (b *Int64) DecodeMsg(r *msgp.Reader) error {
	value, err := r.ReadInt64()
	if err != nil {
		return err
	}
	return b.Scan(value)
}

func (b *Int64) MarshalMsg(i []byte) ([]byte, error) {
	o := msgp.Require(i, b.Msgsize())
	o = msgp.AppendInt64(o, b.Int64)
	return o, nil
}

func (b *Int64) UnmarshalMsg(i []byte) ([]byte, error) {
	value, o, e := msgp.ReadInt64Bytes(i)
	if e != nil {
		return o, e
	}
	return o, b.Scan(value)
}

func (b *Int64) Msgsize() int {
	return msgp.Int64Size
}

func (b *Int64) EncodeMsg(w *msgp.Writer) error {
	return w.WriteInt64(b.Int64)
}

func (nb *NullInt64) EncodeMsg(w *msgp.Writer) error {
	if nb.Valid {
		return w.WriteInt64(nb.Int64)
	}
	return w.WriteNil()
}

func (nb *NullInt64) DecodeMsg(r *msgp.Reader) error {
	if r.IsNil() {
		err := r.ReadNil()
		if err != nil {
			return err
		}
		return nb.Scan(nil)
	}
	value, err := r.ReadInt64()
	if err != nil {
		return err
	}
	return nb.Scan(value)

}

func (nb *NullInt64) MarshalMsg(i []byte) ([]byte, error) {
	if !nb.Valid {
		o := msgp.Require(i, msgp.NilSize)
		o = msgp.AppendNil(o)
		return o, nil
	}
	o := msgp.Require(i, nb.Msgsize())
	o = msgp.AppendInt64(o, nb.Int64)
	return o, nil

}

func (nb *NullInt64) UnmarshalMsg(i []byte) ([]byte, error) {
	t := msgp.NextType(i)
	if t == msgp.NilType {
		o, err := msgp.ReadNilBytes(i)
		if err != nil {
			return o, err
		}
		return o, nb.Scan(nil)
	}
	value, o, err := msgp.ReadInt64Bytes(i)
	if err != nil {
		return o, err
	}
	return o, nb.Scan(value)
}

func (nb *NullInt64) Msgsize() int {
	if !nb.Valid {
		return msgp.NilSize
	}
	return msgp.Int64Size
}
