package channels

import "reflect"

func Drain(c interface{}, recvCallback ...func(v interface{})) {
	cv := reflect.ValueOf(c)
	for {
		x, ok := cv.TryRecv();
		if !ok {
			return
		}
		if len(recvCallback) > 0 {
			recvCallback[0](x.Interface())
		}
	}
}
