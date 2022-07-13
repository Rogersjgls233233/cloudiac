package apps

import (
	"cloudiac/portal/consts/e"
	"cloudiac/portal/libs/ctx"
	"cloudiac/portal/models/forms"
)

// PlatformStatPg 合规策略组数量
func PlatformStatPg(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatPolicy 合规策略数量
func PlatformStatPolicy(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatPgStackEnabled 开启合规并绑定策略组的 Stack 数量
func PlatformStatPgStackEnabled(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatPgEnvEnabledActivate 开启合规并绑定策略组的活跃环境数量
func PlatformStatPgEnvEnabledActivate(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}

// PlatformStatPStackNG 合规不通过的 Stack 数量
func PlatformStatPStackNG(c *ctx.ServiceContext, form *forms.PfStatForm) (interface{}, e.Error) {

	return nil, nil
}
