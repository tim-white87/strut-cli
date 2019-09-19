package provider

import "fmt"

type awsModel struct{}

// NewAwsModel AWS Model constructor
func NewAwsModel() Model {
	return &awsModel{}
}

func (m *awsModel) Provision(r *Resource) {
	fmt.Println(r)
}

func (m *awsModel) Destroy(r *Resource) {

}

func (m *awsModel) CheckStatus() {

}
