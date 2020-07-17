package channels

import "reflect"

func SendNonBlock(c interface{}, v interface{}) {
	cv := reflect.ValueOf(c)
	vv := reflect.ValueOf(v)
	_ = cv.TrySend(vv)
}
