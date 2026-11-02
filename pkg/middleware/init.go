package middleware

import (
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/attachmentPg"
	//_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware/dbPg/postgresqlPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheDomainPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/runnerPg"
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/staterEchoV5"
)
