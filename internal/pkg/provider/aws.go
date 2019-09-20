package provider

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/cecotw/strut-cli/internal/pkg/file"
	"github.com/fatih/color"
)

type awsModel struct {
	session   *session.Session
	resource  *Resource
	template  string
	cfService *cloudformation.CloudFormation
	stack     *cloudformation.Stack
}

// NewAwsModel AWS Model constructor
func NewAwsModel(r *Resource) Model {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	resourceData, err := file.ReadFilePath(r.Path)
	if err != nil {
		color.Red("Issue reading resrouce file path. Path: %s", r.Path)
	}
	template := string(resourceData)
	return &awsModel{
		session:   session,
		resource:  r,
		template:  template,
		cfService: cloudformation.New(session),
	}
}

var iamCapability = "CAPABILITY_NAMED_IAM"
var cababilities = []*string{&iamCapability}

func (m *awsModel) Provision() {
	switch m.CheckStatus() {
	case Status.NotFound:
		color.Green("Creating >>> Resource: %s on Provider: %s", m.resource.Name, m.resource.Provider.Name)
		m.cfService.CreateStack(&cloudformation.CreateStackInput{
			StackName:    &m.resource.Name,
			Capabilities: cababilities,
			TemplateBody: &m.template,
		})
	case Status.Complete:
		color.Green("Updating >>> Resource: %s on Provider: %s", m.resource.Name, m.resource.Provider.Name)
		m.cfService.UpdateStack(&cloudformation.UpdateStackInput{
			StackName:    &m.resource.Name,
			Capabilities: cababilities,
			TemplateBody: &m.template,
		})
	case Status.InProgress:
		color.Yellow("Changes in progress")
	case Status.Failed:
		color.Red("Rolling back >>> Resource: %s on Provider: %s", m.resource.Name, m.resource.Provider.Name)
	}
	m.monitorStackResourcesStatus()
}

func (m *awsModel) Destroy() {
	switch m.CheckStatus() {
	case Status.NotFound:

	case Status.Complete:
		m.cfService.DeleteStack(&cloudformation.DeleteStackInput{
			StackName: m.stack.StackName,
		})
	case Status.InProgress:

	case Status.Failed:
	}
}

func (m *awsModel) CheckStatus() string {
	m.getStacks()
	if m.stack == nil {
		return Status.NotFound
	}
	if *m.stack.StackStatus == cloudformation.StackStatusCreateComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteComplete {
		color.Green("Status >>> Resource: %s has status: %s", *m.stack.StackName, *m.stack.StackStatus)
		return Status.Complete
	} else if *m.stack.StackStatus == cloudformation.StackStatusCreateInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateCompleteCleanupInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusReviewInProgress {
		color.Yellow("Status >>> Resource: %s has status: %s", *m.stack.StackName, *m.stack.StackStatus)
		stackResources, err := m.cfService.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{
			StackName: m.stack.StackName,
		})
		if err != nil {
			color.Red(err.Error())
		}
		for _, stackResource := range stackResources.StackResources {
			color.HiBlack("Resource >>> Status: %s", *stackResource.ResourceStatus)
		}

		return Status.InProgress
	} else if *m.stack.StackStatus == cloudformation.StackStatusCreateFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackFailed {
		color.Red("Status >>> Resource: %s has status: %s", *m.stack.StackName, *m.stack.StackStatus)
	}
	return Status.Failed
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

func (m *awsModel) monitorStackResourcesStatus() {
loop:
	for {
		time.Sleep(time.Second)
		switch m.CheckStatus() {
		case Status.NotFound:
			break loop
		case Status.Complete:
			break loop
		case Status.InProgress:

		case Status.Failed:
			break loop
		}
	}
}
