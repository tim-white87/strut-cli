package provider

type awsModel struct{}

// NewAwsModel AWS Model constructor
func NewAwsModel() Model {
	return &awsModel{}
}

func (m *awsModel) Provision(resources []*Resource) {

}

func (m *awsModel) Destroy() {

}

func (m *awsModel) CheckStatus() {

}
