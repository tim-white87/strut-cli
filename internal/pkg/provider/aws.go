package provider

type awsModel struct{}

// NewAwsModel AWS Model constructor
func NewAwsModel() Model {
	return &awsModel{}
}

func (m *awsModel) Load(p *Provider) {

}

func (m *awsModel) Provision() {

}

func (m *awsModel) Destroy() {

}

func (m *awsModel) CheckStatus() {

}
