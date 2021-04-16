package remote

import (
	"../../../../../tools-go/env"
	"../../../../../tools-go/container"
	"../../../../../tools-go/runtime"
	acc "../../accelerator"
	"../../base"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Prepare() string {
	return acc.F.Prepare(acc.F{}, context, base.F.GetResType(base.F{}, context))
}

func (f F) Create() string {
	return acc.F.Create(acc.F{}, context, base.F.GetNamespace(base.F{}), base.F.GetResType(base.F{}, context), container.F.GetImgAddrs(container.F{}, context, false, false), base.F.GetResources(base.F{}, context), base.F.GetVolumes(base.F{}, context))
}

func (f F) update() {
	acc.F.Update(acc.F{}, context, base.F.GetNamespace(base.F{}), base.F.GetResType(base.F{}, context), container.F.GetImgAddrs(container.F{}, context, false, false), base.F.GetResources(base.F{}, context), base.F.GetVolumes(base.F{}, context))
}

func (f F) Delete() string {
	return acc.F.Delete(acc.F{}, context, base.F.GetNamespace(base.F{}), base.F.GetResType(base.F{}, context), base.F.GetVolumes(base.F{}, context))
}
