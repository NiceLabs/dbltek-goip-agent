package udpstack

//noinspection SpellCheckingInspection
type (
	BCCHUpdate struct {
		RequestID string   `field:"BCCH"`
		ID        string   `field:"id"`
		Password  string   `field:"password"`
		Cells     []string `field:"bcch"`
	}
	CellsUpdate struct {
		RequestID string   `field:"CELLS"`
		ID        string   `field:"id"`
		Password  string   `field:"password"`
		Cells     []string `field:"lists"`
	}
	ATUpdate struct {
		RequestID string `field:"AT"`
		ID        string `field:"id"`
		Password  string `field:"password"`
		Count     uint64 `field:"count"`
		Receive   string `field:"receive"`
	}
)
