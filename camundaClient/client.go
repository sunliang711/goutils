package camundaClient

import (
	"context"
	"fmt"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

type CamundaClient struct {
	client zbc.Client
}

func New(gateway string) (*CamundaClient, error) {
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: gateway}
	client, err := zbc.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("connect to camunda server error: %w", err)
	}

	return &CamundaClient{
		client: client,
	}, nil
}

func (cli *CamundaClient) Close() error {
	return cli.client.Close()
}

func (cli *CamundaClient) DeployProcess(ctx context.Context, name string, processDefinition []byte) (*pb.ProcessMetadata, error) {
	command := cli.client.NewDeployResourceCommand().AddResource(processDefinition, name)

	resource, err := command.Send(ctx)
	if err != nil {
		return nil, err
	}

	if len(resource.GetDeployments()) == 0 {
		return nil, fmt.Errorf("failed to deploy process: %v", name)
	}

	demplyment := resource.GetDeployments()[0]
	process := demplyment.GetProcess()
	if process == nil {
		return nil, fmt.Errorf("failed to deploy process: %v, the deployment was successfule, but no process was returned", name)
	}

	return process, nil
}

func (cli *CamundaClient) StartProcessInstance(ctx context.Context, processId string, vars map[string]any) (*pb.CreateProcessInstanceResponse, error) {
	command, err := cli.client.NewCreateInstanceCommand().BPMNProcessId(processId).LatestVersion().VariablesFromMap(vars)
	if err != nil {
		return nil, err
	}

	process, err := command.Send(ctx)
	if err != nil {
		return nil, err
	}

	return process, nil
}

func (cli *CamundaClient) StartWorker(jobType, workerName string, jobHandler worker.JobHandler, concurrency, maxJobsActive int, timeout, pollInterval time.Duration) worker.JobWorker {
	worker := cli.client.NewJobWorker().JobType(jobType).Handler(jobHandler).Concurrency(concurrency).MaxJobsActive(maxJobsActive).RequestTimeout(timeout).PollInterval(pollInterval).Name(workerName).Open()
	return worker
}
