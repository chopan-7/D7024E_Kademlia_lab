package labCode

/*
NetworkResponse create response messages for network.
All of the RPC are sent using these message functions. The functions will create a Response
object with the data it has to send. Then the MessageHandler function will send and retreive the response
from the other contact.
*/

// Will create a simple ping RPC response object
func (network *Network) CreatePingResponse(res Response) Response {
	responseMessage := Response{
		RPC:            "ping",
		ID:             res.ID,
		SendingContact: &network.Node.Me,
	}
	return responseMessage
}

// Creates a find_node RPC response containing the nodes k closest contacts to a given ID
func (network *Network) CreateFindNodeResponse(res Response) Response {

	contacts := network.Node.Routingtable.FindClosestContacts(res.Body.KadID, bucketSize)

	resBody := Msgbody{
		Nodes: contacts,
	}

	responseMessage := Response{
		RPC:            "find_node",
		ID:             res.ID,
		SendingContact: &network.Node.Me,
		Body:           resBody,
	}

	return responseMessage
}

// Creates a find_data RPC response containing only the data requested if it is stored in the node
// Or will return the 20 closest contacts to the hashed value ID
func (network *Network) CreateFindDataResponse(res Response) Response {
	// Gets the data object from the map if the hash matches a key
	var value []byte
	var containsHash bool
	value, containsHash = network.Node.getDataFromStore(res.Body.Hash)

	if containsHash {

		resBody := Msgbody{
			Data: value,
		}
		responseMessage := Response{
			RPC:            "find_data",
			SendingContact: &network.Node.Me,
			ID:             res.ID,
			Body:           resBody,
		}
		return responseMessage
	}

	contacts := network.Node.Routingtable.FindClosestContacts(NewKademliaID(res.Body.Hash), 20)

	resBody := Msgbody{
		Nodes: contacts,
	}

	responseMessage := Response{
		RPC:            "find_data",
		SendingContact: &network.Node.Me,
		ID:             res.ID,
		Body:           resBody,
	}

	return responseMessage

}

// Creates a simple store_data RPC response to confirm that the data has been stored on the node
func (network *Network) CreateStoreResponse(res Response) Response {
	//Stores data in the node and set the expiration time (TTL)
	network.Node.storeData(res.Body.Data)

	responseMessage := Response{
		RPC:            "store_data",
		SendingContact: &network.Node.Me,
		ID:             res.ID,
	}

	return responseMessage
}
