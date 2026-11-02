package middleware

import (
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg"
	//_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/dbPg/postgresqlPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDomainPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/runnerPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/staterEchoV5"
)
