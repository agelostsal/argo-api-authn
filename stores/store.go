package stores

type Store interface {
	SetUp()
	Close()
	Clone() Store
	QueryServiceTypes(name string) ([]QServiceType, error)
	QueryServiceTypesByUUID(uuid string) ([]QServiceType, error)
	QueryApiKeyAuthMethods(serviceUUID string, host string) ([]QApiKeyAuthMethod, error)
	QueryHeadersAuthMethods(serviceUUID string, host string) ([]QHeadersAuthMethod, error)
	QueryBindingsByAuthID(authID string, serviceUUID string, host string, authType string) ([]QBinding, error)
	QueryBindingsByUUIDAndName(uuid, name string) ([]QBinding, error)
	QueryBindings(serviceUUID string, host string) ([]QBinding, error)
	InsertServiceType(name string, hosts []string, authTypes []string, authMethod string, uuid string, createdOn string, sType string) (QServiceType, error)
	DeleteServiceTypeByUUID(uuid string) error
	InsertAuthMethod(am QAuthMethod) error
	DeleteAuthMethod(am QAuthMethod) error
	DeleteAuthMethodByServiceUUID(serviceUUID string) error
	InsertBinding(name string, serviceUUID string, host string, uuid string, authID string, uniqueKey string, authType string) (QBinding, error)
	UpdateBinding(original QBinding, updated QBinding) (QBinding, error)
	UpdateServiceType(original QServiceType, updated QServiceType) (QServiceType, error)
	UpdateAuthMethod(original QAuthMethod, updated QAuthMethod) (QAuthMethod, error)
	DeleteBinding(qBinding QBinding) error
	DeleteBindingByServiceUUID(serviceUUID string) error
}
