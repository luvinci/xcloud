package common

// 存储类型(表示文件存到哪里)
type StoreType int

const (
	// StoreLocal: 本地存储
	StoreLocal StoreType = 1
	// StoreCeph: Ceph集群
	StoreCeph StoreType = 2
	// StoreOSS: 阿里云OSS
	StoreOSS StoreType = 3
	// StoreMix: 混合(Ceph及OSS)
	StoreMix StoreType = 4
	// StoreAll: 所有类型的存储都存一份数据
	StoreAll StoreType = 5
)
