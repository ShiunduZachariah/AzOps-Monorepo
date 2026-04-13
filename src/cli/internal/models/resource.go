package models

type ResourceGroup struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags,omitempty"`
}

type VirtualMachine struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup,omitempty"`
	VMSize        string            `json:"vmSize,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}
