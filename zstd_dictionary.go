package zstd

/*
#define ZSTD_STATIC_LINKING_ONLY
#define ZBUFF_DISABLE_DEPRECATE_WARNINGS
#include "zstd.h"
#include "zbuff.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func CompressUsingCDict(level int, dict []byte, payload []byte) (res []byte, err error) {
	// create CDict
	cdict := C.ZSTD_createCDict(unsafe.Pointer(&dict[0]), C.size_t(len(dict)), C.int(level))
	if cdict == nil {
		return nil, errors.New("ZSTD_createCDict error")
	}

	res = make([]byte, cCompressBound(len(payload)))

	cctx := C.ZSTD_createCCtx()
	if cctx == nil {
		C.ZSTD_freeCDict(cdict)
		return nil, errors.New("ZSTD_createCCtx() error")
	}

	c_size := C.ZSTD_compress_usingCDict(cctx, unsafe.Pointer(&res[0]), C.size_t(len(res)), unsafe.Pointer(&payload[0]), C.size_t(len(payload)), cdict)

	C.ZSTD_freeCDict(cdict)
	C.ZSTD_freeCCtx(cctx)

	if C.ZSTD_isError(c_size) == 1 {
		return nil, errors.New("error compressing " + C.GoString(C.ZSTD_getErrorName(c_size)))
	}

	return res[:c_size], nil
}
