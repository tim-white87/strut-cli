package provider

import (
	"fmt"

	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/fatih/color"
)

type awsModel struct{}

// NewAwsModel AWS Model constructor
func NewAwsModel() Model {
	return &awsModel{}
}

func (m *awsModel) Provision(r *Resource) {
	resourceData, err := file.ReadFilePath(r.Path)
	if err != nil {
		color.Red("Issue reading resrouce file path. Path: %s", r.Path)
	}
	fmt.Println(resourceData)
}

func (m *awsModel) Destroy(r *Resource) {

}

func (m *awsModel) CheckStatus() {

}
