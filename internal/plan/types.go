package plan

// Plan ist alles, was sich aus einem Manifest und der Hardwareliste ableiten
// laesst, ohne dass der Nutzer es noch einmal angeben muesste.
type Plan struct {
	Roles  []NodeRole
	Ports  []Port
	VLANs  []VLAN
	Links  []Link
	Memory []Memory
}

// NodeRole ordnet einem Knoten aus dem Manifest die Hardware zu, die ihn
// tragen soll.
type NodeRole struct {
	Node     string
	Role     string
	Hardware string
	Vendor   string
	Model    string
	Status   string
	Note     string
}

// Port ist eine Zeile aus der Portbelegung eines Switch-Knotens.
type Port struct {
	Switch     string
	N          int
	Kind       string
	Mode       string
	VLAN       string
	TrunkVLANs []string
	Peer       string
	PoE        bool
	Label      string
}

// VLAN ist ein Eintrag aus dem VLAN-Plan, ergaenzt um die Knoten und Ports,
// die ihn tatsaechlich benutzen.
type VLAN struct {
	Name    string
	ID      int
	Subnet  string
	Egress  string
	Comment string
	UsedBy  []string
}

// Link ist die Streckenpruefung fuer eine physische Verbindung aus
// manifest.links: die kleinere der beiden Portgeschwindigkeiten gegen die
// geforderte Geschwindigkeit.
type Link struct {
	A, B        string
	SpeedGbe    float64
	RequiredGbe float64
	Doubled     bool
	Unknown     bool
	Bottleneck  bool
	Note        string
}

// Memory ist der Arbeitsspeicherbedarf eines Compute-Knotens gegen die
// Dienste, die er tragen soll.
type Memory struct {
	Node       string
	CapacityMB int
	RequiredMB int
	Bottleneck bool
	Note       string
}
