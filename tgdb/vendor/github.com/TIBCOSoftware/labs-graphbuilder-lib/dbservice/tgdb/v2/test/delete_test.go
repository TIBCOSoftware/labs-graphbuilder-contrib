package test

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/connection"
	//"github.com/TIBCOSoftware/tgdb-client/client/goAPI/model"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/query"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	//"github.com/TIBCOSoftware/tgdb-client/client/goAPI/utils"
)

const (
	multiTxnUrl = "tcp://scott@localhost:8222"
	multiTxnPwd = "scott"
)

func CreateToDelete(t *testing.T) {
	fmt.Println("Entering TestCreateToDelete")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(multiTxnUrl, "", multiTxnPwd, nil)
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	err = conn.Connect()
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	gmd, err := conn.GetGraphMetadata(true)
	if err != nil {
		t.Errorf("Returning from SimpleConnectAndValidateBootstrappedEntities - error during conn.GetGraphMetadata")
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	if gof == nil {
		t.Errorf("Returning from MultiTransactionTest - Graph Object Factory is null")
		return
	}

	memberNodeType, err := gmd.GetNodeType("houseMemberType")
	if err != nil {
		t.Errorf("Returning from insertTransaction - error during gmd.GetNodeType('testnode')")
		return
	}
	if memberNodeType != nil {
		t.Logf("'testnode' is found with %d attributes!!", len(memberNodeType.GetAttributeDescriptors()))
	}
	// From Node
	fromNode, err := gof.CreateNodeInGraph(memberNodeType)
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}
	_ = fromNode.SetOrCreateAttribute("memberName", "Napoleon Bonaparte")

	// To Node
	toNode, err := gof.CreateNodeInGraph(memberNodeType)
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}
	_ = toNode.SetOrCreateAttribute("memberName", "Marie Louise of Austria")

	// Edge # 1
	relationEdgeType, err := gmd.GetEdgeType("relation")
	if nil != err {
		t.Errorf("got error : %v", err)
		return
	}

	relationEdge, err := gof.CreateEdgeWithEdgeType(fromNode, toNode, relationEdgeType)
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	_ = relationEdge.SetOrCreateAttribute("relType", "spouse")

	err = conn.DeleteEntity(relationEdge)
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	fmt.Printf("Connection : %v", conn)

	/*
		_, err = conn.Commit()
		if err != nil {
			t.Errorf("got error : %v", err)
			return
		}
	*/
}

func TestSearchToDelete(t *testing.T) {
	fmt.Println("Entering TestSearchToDelete")
	connFactory := connection.NewTGConnectionFactory()
	conn, err := connFactory.CreateConnection(multiTxnUrl, "", multiTxnPwd, nil)
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	err = conn.Connect()
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	gof, err := conn.GetGraphObjectFactory()
	if err != nil {
		t.Errorf("got error : %v", err)
		return
	}

	if gof == nil {
		t.Errorf("Returning from MultiTransactionTest - Graph Object Factory is null")
		return
	}

	// - ======================== query ========================
	// - * queryString           = @nodetype = 'houseMemberType'  and memberName = 'Napoleon Bonaparte';
	// - * edgeFilter            =
	// - * traversalCondition    = @edgetype = 'relation'  and @tonodetype = 'houseMemberType'  and @tonode.memberName = 'Marie Louise of Austria' and @isfromedge = 1 and @degree = 1;
	// - * endCondition          =
	// - * Option.prefetchSize   = 500
	// - * Option.edgeLimit      = 500
	// - * Option.traversalDepth = 1
	// - -------------------------------------------------------

	queryString := "@nodetype = 'houseMemberType'  and memberName = 'Napoleon Bonaparte';"
	traversalCondition := "@edgetype = 'relation'  and @tonodetype = 'houseMemberType'  and @tonode.memberName = 'Marie Louise of Austria' and @isfromedge = 1 and @degree = 1;"

	option := query.NewQueryOption()
	option.SetPreFetchSize(500)
	option.SetTraversalDepth(1)
	option.SetEdgeLimit(500)

	resultSet, err := conn.ExecuteQueryWithFilter(
		queryString,
		"",
		traversalCondition,
		"",
		option,
	)

	if nil == err {
		if nil != resultSet {
			result := resultSet.Next()
			if nil != result {
				fromNode := result.(types.TGNode)
				edges := fromNode.GetEdges()
				for _, edge := range edges {
					toNode := edge.GetVertices()[1]
					if nil != toNode {
						err = conn.DeleteEntity(edge)
						if err != nil {
							t.Errorf("got error : %v", err)
							return
						}

						//_, err = conn.Commit()
						//if err != nil {
						//	t.Errorf("got error : %v", err)
						//	return
						//}
					}
				}
			}
		}
	}
}
