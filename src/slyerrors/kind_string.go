// Code generated by "stringer -type=Kind -linecomment"; DO NOT EDIT.

package slyerrors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[KindUnknown-0]
	_ = x[KindUnexpected-1]
	_ = x[KindUnauthorized-2]
	_ = x[KindForbidden-3]
	_ = x[KindValidation-4]
	_ = x[KindBadRequest-5]
	_ = x[KindNotFound-6]
	_ = x[KindConflict-7]
	_ = x[KindUnprocessableEntity-8]
}

const _Kind_name = "UnknownUnexpectedUnauthorizedForbiddenValidationBadRequestNotFoundConflictUnprocessabeEntity"

var _Kind_index = [...]uint8{0, 7, 17, 29, 38, 48, 58, 66, 74, 92}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
