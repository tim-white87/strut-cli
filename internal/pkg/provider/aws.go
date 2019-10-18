package provider

import (
	"sync"
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
		_, err := m.cfService.CreateStack(&cloudformation.CreateStackInput{
			StackName:    &m.resource.Name,
			Capabilities: cababilities,
			TemplateBody: &m.template,
		})
		if err != nil {
			color.Red("Create Error >>> %s: %s", m.resource.Name, err.Error())
		}
	case Status.Complete:
		color.Green("Updating >>> Resource: %s on Provider: %s", m.resource.Name, m.resource.Provider.Name)
		_, err := m.cfService.UpdateStack(&cloudformation.UpdateStackInput{
			StackName:    &m.resource.Name,
			Capabilities: cababilities,
			TemplateBody: &m.template,
		})
		if err != nil {
			color.Red("Update Error >>> %s: %s", m.resource.Name, err.Error())
		}
	case Status.InProgress:
		color.Yellow("Changes in progress")
	case Status.Failed:
		color.Red("Rolling back >>> Resource: %s on Provider: %s", m.resource.Name, m.resource.Provider.Name)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go m.monitorStackResourcesStatus(wg)
}

func (m *awsModel) Destroy() {
	// TODO empty S3 buckets
	// TODO Delete Origin Access Identity
	// TODO delete all cognito users
	status := m.CheckStatus()
	if status == Status.Complete || status == Status.Failed {
		color.Green("Deleting >>> Resource: %s on Provider: %s", m.resource.Name, m.resource.Provider.Name)
		_, err := m.cfService.DeleteStack(&cloudformation.DeleteStackInput{
			StackName: m.stack.StackName,
		})
		if err != nil {
			color.Red("Destroy Error >>> %s: %s", m.resource.Name, err.Error())
		}
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go m.monitorStackResourcesStatus(wg)
}

func (m *awsModel) CheckStatus() string {
	var lastStatus string
	if m.stack != nil {
		lastStatus = *m.stack.StackStatus
	}
	m.getStack()
	if m.stack == nil {
		return Status.NotFound
	}
	var status string
	if *m.stack.StackStatus == cloudformation.StackStatusCreateComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteComplete ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackComplete {
		status = Status.Complete
	} else if *m.stack.StackStatus == cloudformation.StackStatusCreateInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateCompleteCleanupInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress ||
		*m.stack.StackStatus == cloudformation.StackStatusReviewInProgress {
		status = Status.InProgress
	} else if *m.stack.StackStatus == cloudformation.StackStatusCreateFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusRollbackFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusDeleteFailed ||
		*m.stack.StackStatus == cloudformation.StackStatusUpdateRollbackFailed {
		status = Status.Failed
	}
	if lastStatus != *m.stack.StackStatus {
		color.White("Resource: %s >>> Status: %s", *m.stack.StackName, *m.stack.StackStatus)
		if m.stack.StackStatusReason != nil {
			color.HiBlack("Reason: %s", *m.stack.StackStatusReason)
		}
	}
	return status
}

func (m *awsModel) getStack() {
	stackOutput, _ := m.cfService.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: &m.resource.Name,
	})
	if len(stackOutput.Stacks) > 0 {
		m.stack = stackOutput.Stacks[0]
	}
}

func (m *awsModel) monitorStackResourcesStatus(wg *sync.WaitGroup) {
	defer wg.Done()
	var status string
loop:
	for {
		time.Sleep(time.Second)
		status = m.CheckStatus()
		if status != Status.InProgress {
			break loop
		}
	}
}
