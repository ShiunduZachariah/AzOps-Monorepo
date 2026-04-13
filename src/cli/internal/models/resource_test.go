package models

import (
	"encoding/json"
	"testing"
)

func TestResourceGroupMarshalJSON(t *testing.T) {
	t.Parallel()

	group := ResourceGroup{
		ID:       "id-1",
		Name:     "rg-one",
		Location: "eastus",
	}

	data, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("marshal resource group: %v", err)
	}

	got := string(data)
	if got != `{"id":"id-1","name":"rg-one","location":"eastus"}` {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestVirtualMachineMarshalJSON(t *testing.T) {
	t.Parallel()

	vm := VirtualMachine{
		ID:            "vm-1",
		Name:          "vm-one",
		Location:      "eastus",
		ResourceGroup: "rg-one",
		VMSize:        "Standard_B2s",
	}

	data, err := json.Marshal(vm)
	if err != nil {
		t.Fatalf("marshal virtual machine: %v", err)
	}

	got := string(data)
	if got != `{"id":"vm-1","name":"vm-one","location":"eastus","resourceGroup":"rg-one","vmSize":"Standard_B2s"}` {
		t.Fatalf("unexpected json: %s", got)
	}
}
