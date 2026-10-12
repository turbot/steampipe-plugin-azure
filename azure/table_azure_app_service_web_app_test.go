package azure

import "testing"

func TestIsFunctionAppKind(t *testing.T) {
	kind := func(s string) *string { return &s }
	cases := []struct {
		kind *string
		want bool
	}{
		{nil, false},
		{kind("app"), false},
		{kind("app,linux"), false},
		{kind("app,linux,container"), false},
		{kind("functionapp"), true},
		{kind("functionapp,linux"), true},
		{kind("functionapp,workflowapp"), true},
		{kind("functionapp,linux,container,workflowapp"), true},
	}
	for _, c := range cases {
		name := "<nil>"
		if c.kind != nil {
			name = *c.kind
		}
		if got := isFunctionAppKind(c.kind); got != c.want {
			t.Errorf("isFunctionAppKind(%q) = %v, want %v", name, got, c.want)
		}
	}
}
