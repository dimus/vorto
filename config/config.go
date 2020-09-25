package config

type Config struct {
	DataDir    string
	DefaultSet string
	WordsBatch int
	Sets       map[string]struct{}
}

func NewConfig(opts ...Option) Config {
	cnf := Config{}
	for _, opt := range opts {
		opt(&cnf)
	}
	return cnf
}

type Option func(*Config)

func OptDataDir(s string) Option {
	return func(cnf *Config) {
		cnf.DataDir = s
	}
}

func OptDefaultSet(s string) Option {
	return func(cnf *Config) {
		cnf.DefaultSet = s
	}
}

func OptWordsBatch(i int) Option {
	return func(cnf *Config) {
		cnf.WordsBatch = i
	}
}

func OptSets(sets []string) Option {
	return func(cnf *Config) {
		res := make(map[string]struct{})
		for _, v := range sets {
			res[v] = struct{}{}
		}
		cnf.Sets = res
	}
}
