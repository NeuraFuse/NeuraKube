package base

import (
	"../../../../tools-go/apps/python/debugpy"
	"../../../../tools-go/apps/python/flask"
	"../../../../tools-go/apps/tensorflow/tensorboard"
	"../../../../tools-go/config"
	"../../../../tools-go/env"
	"../../../../tools-go/errors"
	"../../../../tools-go/kubernetes/namespaces"
	"../../../../tools-go/objects/strings"
	"../../../../tools-go/projects"
	"../../../../tools-go/runtime"
	"../../../../tools-go/users"
)

type F struct{}

func (f F) GetVolumeSizeGB(context string) string {
	return config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".VolumeSizeGB", "") + "Gi"
}

func (f F) GetVolumeContainerPath() string {
	return env.F.GetContainerWorkingDir(env.F{}) + projects.F.GetWorkingDir(projects.F{})
}

func (f F) GetVolumes(context string) [][]string {
	return [][]string{{f.GetVolumeContainerPath(), f.GetVolumeSizeGB(context)}}
}

func (f F) GetVolumeMountPath(envFramework string) string {
	return env.F.GetContainerWorkingDir(env.F{}) + "/lightning/" + envFramework
}

func (f F) GetNamespace() string {
	return namespaces.Default + "-" + users.GetIDActive()
}

func (f F) GetResType(context string) string {
	return f.GetEnvFramework(context, false) + " " + context + " environment [" + f.GetType(context, false) + "]"
}

func (f F) GetEnvFramework(context string, titleCase bool) string {
	var envFramework string = config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Environment.Framework", "")
	if titleCase {
		envFramework = strings.Title(envFramework)
	}
	return envFramework
}

func (f F) GetResources(context string) string {
	return config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Type", "")
}

func (f F) GetContainerPorts(context string) [][]string {
	var apps []string
	var appsBase []string = []string{"debugpy", "tensorboard"}
	switch context {
	case "develop":
		apps = appsBase
	case "app":
		apps = appsBase
	case "inference":
		apps = []string{"flask"}
	}
	return f.GetContainerPortsForApps(apps)
}

func (f F) GetContainerPortsForApps(apps []string) [][]string {
	var containerPorts [][]string
	for _, app := range apps {
		switch app {
		case "debugpy":
			{
				containerPorts = append(containerPorts, debugpy.ContainerPorts)
			}
		case "tensorboard":
			{
				containerPorts = append(containerPorts, tensorboard.ContainerPorts)
			}
		case "flask":
			{
				containerPorts = append(containerPorts, flask.ContainerPorts)
			}
		default:
			{
				errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There are no registred ports for the "+app+"!", true, true, true)
			}
		}
	}
	return containerPorts
}

func (f F) GetType(context string, upperCase bool) string {
	if context == "develop/remote" { // TODO: Refactor
		context = "remote"
	}
	var resType string = config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Type", "")
	if upperCase {
		resType = strings.ToUpper(resType)
	}
	return resType
}
