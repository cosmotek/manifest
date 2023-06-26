package scanner

type Asset struct {
	Identifier string
	Metadata   map[string]any
}

type AssetList []Asset

func (r AssetList) Len() int {
	return len(r)
}
func (r AssetList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r AssetList) Less(i, j int) bool {
	return len(r[i].Identifier) < len(r[j].Identifier)
}

type EnvironmentScanner interface {
	RunScan() (AssetList, error)
}
