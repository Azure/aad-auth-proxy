package contracts

//
// Contract to obtain access token for requests to ingestion endpoint.
//
type ITokenProvider interface {
	GetAccessToken() (string, error)
}
