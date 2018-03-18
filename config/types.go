package config

type Configuration struct {
	Repositories []repository `toml:"repository"`
}

type repository struct {
	Name    string   `toml:"name"`
	Mirrors []Mirror `toml:"mirror"`
}

type Mirror struct {
	Name       string `toml:"name"`
	Remote     string `toml:"remote"`
	AuthMode   string `toml:"auth_mode"`
	SSHKey     string `toml:"ssh_key"`
	Passphrase string `toml:"passphrase"`
	User string `toml:"user"`
	Password string `toml:"password"`
}
