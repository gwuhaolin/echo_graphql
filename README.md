[echo](https://echo.labstack.com/) with [graphql-go](https://github.com/graphql-go/graphql)

### Basic use
```go
import (
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
)

e := echo.New()
e.Any("/graphql", echo_graphql.NewEchoHandle(echo_graphql.EchoHandleOptions{
	Schema: graphql.MustParseSchema(`your graphql schema define content...`),
}))
```

### Use cache to improve performance  
cache all graphql request by post body hash
```go
import (
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
	"github.com/gwuhaolin/lfucache"
)

e := echo.New()
graphqlSchema := graphql.MustParseSchema(`your graphql schema define content...`)

e.Any("/graphql", echo_graphql.NewEchoHandle(echo_graphql.EchoHandleOptions{
	Schema: graphqlSchema,
    Cache:  lfucache.NewLfuCache(1024),
}))
```

skip cache some request 
```go
import (
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
	"github.com/gwuhaolin/lfucache"
)

e := echo.New()
graphqlSchema := graphql.MustParseSchema(`your graphql schema define content...`)

e.Any("/graphql", echo_graphql.NewEchoHandle(echo_graphql.EchoHandleOptions{
	Schema: graphqlSchema,
    Cache:  lfucache.NewLfuCache(1024),
    SkipCache: func(params *Params){
        retuen true // don't cache this request
    }
}))
```