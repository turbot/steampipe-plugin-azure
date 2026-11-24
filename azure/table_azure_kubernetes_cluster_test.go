package azure

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
)

type fakeClustersPager struct {
	pages [][]*armcontainerservice.ManagedCluster
	errAt int // index of the page that fails, -1 for none
	next  int
}

func (p *fakeClustersPager) More() bool { return p.next < len(p.pages) }

func (p *fakeClustersPager) NextPage(context.Context) (armcontainerservice.ManagedClustersClientListResponse, error) {
	i := p.next
	p.next++
	if i == p.errAt {
		return armcontainerservice.ManagedClustersClientListResponse{}, errors.New("403 AuthorizationFailed")
	}
	var r armcontainerservice.ManagedClustersClientListResponse
	r.Value = p.pages[i]
	return r, nil
}

func clusters(names ...string) []*armcontainerservice.ManagedCluster {
	out := make([]*armcontainerservice.ManagedCluster, 0, len(names))
	for _, n := range names {
		n := n
		out = append(out, &armcontainerservice.ManagedCluster{Name: &n})
	}
	return out
}

func TestStreamManagedClusters(t *testing.T) {
	ctx := context.Background()

	t.Run("all pages", func(t *testing.T) {
		p := &fakeClustersPager{pages: [][]*armcontainerservice.ManagedCluster{clusters("a", "b"), clusters("c")}, errAt: -1}
		var got []string
		err := streamManagedClusters(ctx, p, func(c *armcontainerservice.ManagedCluster) bool { got = append(got, *c.Name); return true })
		if err != nil || len(got) != 3 {
			t.Fatalf("got %v, err %v; want 3 clusters and no error", got, err)
		}
	})

	t.Run("page error is returned", func(t *testing.T) {
		p := &fakeClustersPager{pages: [][]*armcontainerservice.ManagedCluster{clusters("a"), nil}, errAt: 1}
		n := 0
		err := streamManagedClusters(ctx, p, func(*armcontainerservice.ManagedCluster) bool { n++; return true })
		if err == nil {
			t.Fatalf("want the page error, got nil after %d clusters", n)
		}
	})

	t.Run("error on first page is not an empty result", func(t *testing.T) {
		p := &fakeClustersPager{pages: [][]*armcontainerservice.ManagedCluster{nil}, errAt: 0}
		if err := streamManagedClusters(ctx, p, func(*armcontainerservice.ManagedCluster) bool { return true }); err == nil {
			t.Fatal("want an error, got nil")
		}
	})

	t.Run("stops when emit returns false", func(t *testing.T) {
		p := &fakeClustersPager{pages: [][]*armcontainerservice.ManagedCluster{clusters("a", "b"), clusters("c")}, errAt: -1}
		n := 0
		err := streamManagedClusters(ctx, p, func(*armcontainerservice.ManagedCluster) bool { n++; return n < 2 })
		if err != nil || n != 2 {
			t.Fatalf("emitted %d, err %v; want 2 and no error", n, err)
		}
	})
}
