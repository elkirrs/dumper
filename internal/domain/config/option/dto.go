package option

type Options struct {
	AuthSource string   `yaml:"auth_source"`
	SSL        *bool    `yaml:"ssl" default:"false"`
	Mode       string   `yaml:"mode"`
	Role       string   `yaml:"role"`
	Path       string   `yaml:"path"`
	Source     string   `yaml:"source"`
	IncTables  []string `yaml:"inc_tables"`
	ExcTables  []string `yaml:"exc_tables"`
	Host       string   `yaml:"host" default:"http://127.0.0.1"`

	// DynamoDB specific options
	Region   string `yaml:"region"`   // AWS region
	Profile  string `yaml:"profile"`  // AWS profile name
	Endpoint string `yaml:"endpoint"` // Custom endpoint (for local DynamoDB)

	//InfluxDB
	Bucket         string `yaml:"bucket,omitempty"`
	BucketId       string `yaml:"bucket_id,omitempty"`
	Organization   string `yaml:"org,omitempty"`
	OrganizationId string `yaml:"org_id,omitempty"`
	Start          string `yaml:"start"`
	End            string `yaml:"end"`
	Filter         string `yaml:"filter"`
	SkipVerify     *bool  `yaml:"skip_verify" default:"false"`
	NodeId         string `yaml:"node_id"`
	DataDir        string `yaml:"data_dir"`

	Version    string `yaml:"version"`
	BackupMode string `yaml:"backup_mode"`

	//Firebird
	SkipIssue     bool `yaml:"skip_issue" default:"false"`
	FastAndStable bool `yaml:"fast_and_stable" default:"false"`
	SkipGarbage   bool `yaml:"skip_garbage" default:"false"`

	// OpenSearch and ElasticSearch
	CACertPath         string   `yaml:"ca_crt_path"`
	KeyPath            string   `yaml:"key_path"`
	CertPath           string   `yaml:"cert_path"`
	KeyPass            string   `yaml:"key_pass"`
	SnapPath           string   `yaml:"snap_path"`
	Indices            []string `yaml:"indices"`
	IgnoreUnavailable  *bool    `yaml:"ignore_unavailable" default:"false"`
	IncludeGlobalState *bool    `yaml:"include_global_state" default:"false"`
}
