package ssnet

type GetRequest struct {
	URL    string
	SSlist SteppingStones
}

func NewGetRequest(u string, sslist SteppingStones) GetRequest {
	return GetRequest{
		URL:    u,
		SSlist: sslist,
	}
}
