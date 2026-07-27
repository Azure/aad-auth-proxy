package contracts

type IHostTokenProvider interface {
    ITokenProvider
    GetTokenForHost(host string) (string, error)
}