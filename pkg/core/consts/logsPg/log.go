package logsPg

import (
	"go-spring.org/log"
)

var (
	// TagAppDef is the default tag used for application logs.
	TagAppDef = RegisterAppTagPg("app", "")

	// TagBizDef is the default tag used for business-related logs.
	TagBizDef = RegisterBizTagPg("biz", "")
)

// RegisterAppTag returns a Tag used for application-layer logs (e.g., framework events, lifecycle).
// subType represents the component or module, action represents the lifecycle phase or behavior.
func RegisterAppTagPg(subType, action string) *log.Tag {
	return log.RegisterTag(log.BuildTag("pg", subType, action))
}

// RegisterBizTag returns a Tag used for business-logic logs (e.g., use cases, domain events).
// subType is the business domain or feature name, action is the operation being logged.
func RegisterBizTagPg(subType, action string) *log.Tag {
	return log.RegisterTag(log.BuildTag("pg", subType, action))
}
