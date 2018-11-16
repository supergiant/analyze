package cloudprovider

type ComputeInstance struct {
	InstanceID   string `json:"instanceId"`
	InstanceType string `json:"instanceType"`
}

type ProductPrice struct {
	InstanceType string `json:"instanceType"`
	Memory       string `json:"memory"`
	Vcpu         string `json:"vcpu"`
	Unit         string `json:"unit"`
	Currency     string `json:"currency"`
	ValuePerUnit string `json:"valuePerUnit"`
	UsageType    string `json:"usageType"`
	Tenancy      string `json:"tenancy"`
}
