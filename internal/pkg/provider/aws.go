package provider

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/fatih/color"
)

type awsModel struct {
	resource *Resource
	session  *session.Session
}

// NewAwsModel AWS Model constructor
func NewAwsModel() Model {
	return &awsModel{
		session: session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})),
	}
}

func (m *awsModel) Provision(r *Resource) {
	m.resource = r
	resourceData, err := file.ReadFilePath(m.resource.Path)
	if err != nil {
		color.Red("Issue reading resrouce file path. Path: %s", m.resource.Path)
	}
	template := string(resourceData)
	m.deployCloudFormationStack(template)
}

func (m *awsModel) Destroy(r *Resource) {

}

func (m *awsModel) CheckStatus() {

}

func (m *awsModel) deployCloudFormationStack(template string) {
	fmt.Println(m.session)
}
