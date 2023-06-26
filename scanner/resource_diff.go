package scanner

import (
	"github.com/google/go-cmp/cmp"
)

var compareOpts = []cmp.Option{
	// cmpopts.IgnoreUnexported(
	// 	resourceexplorer2.Resource{},
	// ),
	// cmpopts.IgnoreFields(
	// 	resourceexplorer2.Resource{},
	// 	"LastReportedAt",
	// ),
}

func ComputeDiff(oldAssets, newAssets AssetList) *string {
	if len(oldAssets) > 0 && !cmp.Equal(oldAssets, newAssets) {
		v := cmp.Diff(oldAssets, newAssets, compareOpts...)
		return &v
	}

	return nil
}
