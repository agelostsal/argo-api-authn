package stores

type QService struct {
	Name           string   `json:"name" bson:"name"`
	Hosts          []string `json:"hosts" bson:"hosts"`
	AuthTypes      []string `json:"auth_types" bson:"auth_types"`
	AuthMethod     string   `json:"auth_method" bson:"auth_method"`
	RetrievalField string   `json:"retrieval_field" bson:"retrieval_field"`
}

type QBinding struct {
	Name      string `json:"name" bson:"name"`
	Service   string `json:"service" bson:"service"`
	Host      string `json:"host" bson:"host"`
	DN        string `json:"dn,omitempty" bson:"dn,omitempty"`
	OIDCToken string `json:"oidc_token,omitempty" bson:"oidc_token,omitempty"`
	UniqueKey string `json:"unique_key,omitempty" bson:"unique_key,omitempty"`
	CreatedOn string `json:"created_on,omitempty" bson:"created_on,omitempty"`
	LastAuth  string `json:"last_auth,omitempty" bson:"last_auth,omitempty"`
}

type QHost struct {
	Service     string `json:"service" bson:"service"`
	Name        string `json:"name" bson:"name"`
	Port        int    `json:"port" bson:"port"`
	AccessToken string `json:"access_token" bson:"access_token"`
	Path        string `json:"path" bson:"path"`
}

type QApiKeyAuth struct {
	Type      string `json:"type" bson:"type"`
	Service   string `json:"service" bson:"service"`
	Host      string `json:"host" bson:"host"`
	Port      int    `json:"port" bson:"port"`
	Path      string `json:"path" bson:"path"`
	AccessKey string `json:"access_key" bson:"access_key"`
}
