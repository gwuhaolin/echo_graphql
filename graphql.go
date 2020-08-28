package echo_graphql

import (
	"encoding/hex"
	"encoding/json"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"unsafe"

	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
)

type Params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func hashBody(body io.ReadCloser) (string, []byte, error) {
	bs, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return "", nil, err
	}
	hash := crc32.ChecksumIEEE(bs)
	key := hex.EncodeToString((*[4]byte)(unsafe.Pointer(&hash))[:])
	return key, bs, nil
}

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
}

type EchoHandleOptions struct {
	Schema    *graphql.Schema
	Cache     Cache
	SkipCache func(params *Params) bool
}

func NewEchoHandle(options EchoHandleOptions) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		key, bs, err := hashBody(context.Request().Body)
		if err != nil {
			return err
		}
		var ret *graphql.Response
		if v, has := options.Cache.Get(key); has {
			ret = v.(*graphql.Response)
		} else {
			params := new(Params)
			if err = json.Unmarshal(bs, params); err != nil {
				return
			}
			ret = options.Schema.Exec(context.Request().Context(), params.Query, params.OperationName, params.Variables)
			go func() {
				if options.SkipCache != nil && options.SkipCache(params) {
					return
				}
				options.Cache.Set(key, ret)
			}()
		}
		return context.JSON(http.StatusOK, ret)
	}
}
