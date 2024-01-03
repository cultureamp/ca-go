# ca-go/log

words

## Environment Variables

YOU MUST set these:
- APP_NAME
- LOG_LEVEL (default INFO)


## Examples


```
package cago

import (
	"context"

	"github.com/cultureamp/ca-go/log"
	"github.com/cultureamp/ca-go/jwt"
)


func test() {
	ctx := context.Background()

    req := http.NewRequest(....)
    ctx := log.ContextWithRequest(ctx, req)

    JwtPayload := jwt.Decode(AUTH_HEADER) 
    authUsers.UserID := JwtPayLoad.UserID
    authUsers.AccountID := JwtPayLoad.AccountID
    authUsers.RealUserID := JwtPayLoad.ReadUserID
    ctx := log.ContextWithAuthUserIDs(ctx, JwtPayload)

	log.Debug(ctx, "hello")
}
```
