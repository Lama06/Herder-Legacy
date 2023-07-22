package astar

type Neighbour struct {
	Node Node
	Cost float64
}

type Node interface {
	AstarNeighbours(context any) []Neighbour

	AstarEstimateCost(context any, to Node) float64
}

func FindShortestPath(context any, from Node, to Node) (path []Node, found bool) {
	type nodeDaten struct {
		offen                      bool
		parent                     Node
		costFromStart              float64
		extimatedCostToDestination float64
	}
	totalCost := func(node nodeDaten) float64 {
		return node.costFromStart + node.extimatedCostToDestination
	}

	bekannteNodes := map[Node]nodeDaten{
		from: {
			offen:                      true,
			parent:                     nil,
			costFromStart:              0,
			extimatedCostToDestination: from.AstarEstimateCost(context, to),
		},
	}

	for {
		var (
			nächsteNode      Node
			nächsteNodeDaten nodeDaten
		)
		for node, nodeDaten := range bekannteNodes {
			if !nodeDaten.offen {
				continue
			}
			if nächsteNode == nil ||
				totalCost(nodeDaten) < totalCost(nächsteNodeDaten) ||
				(totalCost(nodeDaten) == totalCost(nächsteNodeDaten) && nodeDaten.extimatedCostToDestination < nächsteNodeDaten.extimatedCostToDestination) {
				nächsteNode = node
				nächsteNodeDaten = nodeDaten
			}
		}

		switch {
		case nächsteNode == nil:
			return nil, false
		case nächsteNode == to:
			path := []Node{to}
			for bekannteNodes[path[0]].parent != nil {
				path = append([]Node{bekannteNodes[path[0]].parent}, path...)
			}
			return path, true
		default:
			nächsteNodeDaten.offen = false
			bekannteNodes[nächsteNode] = nächsteNodeDaten

			for _, neighbour := range nächsteNode.AstarNeighbours(context) {
				neighbourDaten, istBekannt := bekannteNodes[neighbour.Node]
				switch {
				case istBekannt && !neighbourDaten.offen:
					continue
				case istBekannt && nächsteNodeDaten.costFromStart+neighbour.Cost < neighbourDaten.costFromStart:
					neighbourDaten.costFromStart = nächsteNodeDaten.costFromStart + neighbour.Cost
					bekannteNodes[neighbour.Node] = neighbourDaten
				case !istBekannt:
					bekannteNodes[neighbour.Node] = nodeDaten{
						offen:                      true,
						parent:                     nächsteNode,
						costFromStart:              nächsteNodeDaten.costFromStart + neighbour.Cost,
						extimatedCostToDestination: neighbour.Node.AstarEstimateCost(context, to),
					}
				}
			}
		}
	}
}
