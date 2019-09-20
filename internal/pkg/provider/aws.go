package provider

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/fatih/color"
)

type awsModel struct {
	resource  *Resource
	session   *session.Session
	cfService *cloudformation.CloudFormation
	stack     *cloudformation.Stack
}

// NewAwsModel AWS Model constructor
func NewAwsModel() Model {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &awsModel{
		session:   session,
		cfService: cloudformation.New(session),
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

func (m *awsModel) CheckStatus() bool {
	m.getStacks()
	if *m.stack.StackStatus == cloudformation.StackStatusCreateComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteComplete {
		color.Green(*m.stack.StackStatus)
		return true
	} else if *m.stack.StackStatus == cloudformation.StackStatusCreateInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateCompleteCleanupInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusReviewInProgress {
		color.Yellow(*m.stack.StackStatus)
	} else if *m.stack.StackStatus == cloudformation.StackStatusCreateFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackFailed {
		color.Red(*m.stack.StackStatus)
	}
	return false
}

func (m *awsModel) deployCloudFormationStack(template string) {
	m.CheckStatus()

	//  resourceExists := false

	// if resourceExists
}

func (m *awsModel) getStacks() {
	stackOutput, err := m.cfService.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: &m.resource.Name,
	})
	if err != nil {
		color.Red(err.Error())
	}
	if len(stackOutput.Stacks) > 0 {
		m.stack = stackOutput.Stacks[0]
	}

}
