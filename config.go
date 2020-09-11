package gowalker

import "os"

func Config(cfg interface{}) error {
	return Walk(cfg,
		Tag("default"),
		DefaultEnvWalker(),
		DefaultFlagWalker(),
	)
}

func DefaultEnvWalker() Walker {
	return StringSetter(
		FieldKeys(Tag("env"), EnvNamer),
		LookupFuncSource(os.LookupEnv),
	)
}
