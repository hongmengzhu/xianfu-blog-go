package staterEchoV5

import "go-spring.org/spring/gs"

func init() {
	gs.Provide(NewServer, gs.IndexArg(1, gs.TagArg("${spring.echo.server}"))).Export(gs.As[gs.Server]()).
		Condition(gs.And(
			gs.OnProperty("spring.echo.server.enabled").HavingValue("true").MatchIfMissing(),
			gs.OnProperty("spring.echo.server.addr"),
		))
}
