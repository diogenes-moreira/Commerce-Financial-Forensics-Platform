package cloudaccount

type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderGCP   Provider = "gcp"
	ProviderAzure Provider = "azure"
)

func (p Provider) IsValid() bool {
	switch p {
	case ProviderAWS, ProviderGCP, ProviderAzure:
		return true
	default:
		return false
	}
}
