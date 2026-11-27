package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func TestNewSlotInfo(t *testing.T) {
	str := func(s string) *string { return &s }
	slot := web.Site{
		ID:       str("/subscriptions/s/resourceGroups/rg/providers/Microsoft.Web/sites/app/slots/staging"),
		Name:     str("app/staging"),
		Kind:     str("functionapp,linux"),
		Location: str("westeurope"),
		Type:     str("Microsoft.Web/sites/slots"),
		Identity: &web.ManagedServiceIdentity{Type: web.ManagedServiceIdentityTypeSystemAssigned, PrincipalID: str("p")},
	}

	info := newSlotInfo(slot, "app")
	if *info.AppName != "app" || info.ID != slot.ID || info.Name != slot.Name || info.Kind != slot.Kind {
		t.Fatalf("newSlotInfo did not copy the slot: %+v", info)
	}

	// Every page of a list and the get stream this type, and the identity transform and the
	// configuration hydrate assert it, so a raw web.Site on any path fails the query
	got, err := slotIdentity(context.Background(), &transform.TransformData{HydrateItem: info})
	if err != nil {
		t.Fatal(err)
	}
	if identity := got.(map[string]interface{}); identity["PrincipalID"] != slot.Identity.PrincipalID {
		t.Errorf("slotIdentity = %v, want the slot's identity", identity)
	}

	// Each call keeps its own app name
	if other := newSlotInfo(slot, "other"); *info.AppName != "app" || *other.AppName != "other" {
		t.Errorf("app names leaked between items: %q, %q", *info.AppName, *other.AppName)
	}
}
