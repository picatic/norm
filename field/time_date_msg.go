package field

import (
	"github.com/tinylib/msgp/msgp"
)

func (b *TimeDate) DecodeMsg(r *msgp.Reader) error {
	value, err := r.ReadTime()
	if err != nil {
		return err
	}
	return b.Scan(value)
}

func (b *TimeDate) MarshalMsg(i []byte) ([]byte, error) {
	o := msgp.Require(i, b.Msgsize())
	o = msgp.AppendTime(o, b.Time)
	return o, nil
}

func (b *TimeDate) UnmarshalMsg(i []byte) ([]byte, error) {
	value, o, e := msgp.ReadTimeBytes(i)
	if e != nil {
		return o, e
	}
	return o, b.Scan(value)
}

func (b *TimeDate) Msgsize() int {
	return msgp.TimeSize
}

func (b *TimeDate) EncodeMsg(w *msgp.Writer) error {
	return w.WriteTime(b.Time)
}

func (nb *NullTimeDate) EncodeMsg(w *msgp.Writer) error {
	if nb.Valid {
		return w.WriteTime(nb.Time)
	}
	return w.WriteNil()
}

func (nb *NullTimeDate) DecodeMsg(r *msgp.Reader) error {
	if r.IsNil() {
		err := r.ReadNil()
		if err != nil {
			return err
		}
		return nb.Scan(nil)
	}
	value, err := r.ReadTime()
	if err != nil {
		return err
	}
	return nb.Scan(value)

}

func (nb *NullTimeDate) MarshalMsg(i []byte) ([]byte, error) {
	if !nb.Valid {
		o := msgp.Require(i, msgp.NilSize)
		o = msgp.AppendNil(o)
		return o, nil
	}
	o := msgp.Require(i, nb.Msgsize())
	o = msgp.AppendTime(o, nb.Time)
	return o, nil

}

func (nb *NullTimeDate) UnmarshalMsg(i []byte) ([]byte, error) {
	t := msgp.NextType(i)
	if t == msgp.NilType {
		o, err := msgp.ReadNilBytes(i)
		if err != nil {
			return o, err
		}
		return o, nb.Scan(nil)
	}
	value, o, err := msgp.ReadTimeBytes(i)
	if err != nil {
		return o, err
	}
	return o, nb.Scan(value)
}

func (nb *NullTimeDate) Msgsize() int {
	if !nb.Valid {
		return msgp.NilSize
	}
	return msgp.TimeSize
}
