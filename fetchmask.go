package main

type FetchMask uint8

const (
	FMonitors FetchMask = 1 << iota
	FWorkspaces
	FClients
	FActiveWS
	FActiveWin
)

func (mask FetchMask) Has(flag FetchMask) bool {
	return mask&flag != 0
}

const (
	// For "workspace goto N"
	MaskGoto = FWorkspaces | FClients | FActiveWS | FMonitors

	// For "workspace move N [--all]"
	MaskMove = FWorkspaces | FClients | FActiveWS | FActiveWin | FMonitors

	// For "workspace cycle up|down"
	MaskCycle = FWorkspaces | FClients | FActiveWS | FMonitors
)
