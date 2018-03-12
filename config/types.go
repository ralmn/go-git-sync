package config

type Configuration struct {
	Repositories []repository `toml:"repository"`
}

type repository struct {
	Name    string   `toml:"name"`
	Mirrors []Mirror `toml:"mirror"`
}

type Mirror struct {
	Name    string `toml:"name"`
	Remote  string `toml:"remote"`
	SSH_KEY string `toml:"ssh_key"`
}
