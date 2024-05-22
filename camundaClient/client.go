package camundaClient

import (
	"context"
	"fmt"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/commands"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"google.golang.org/grpc"
)

type CamundaClient struct {
	Client zbc.Client
}

func New(gateway string) (*CamundaClient, error) {
	config := zbc.ClientConfig{
		UsePlaintextConnection: true,
		GatewayAddress:         gateway,
		DialOpts:               []grpc.DialOption{grpc.WithBlock(), grpc.WithTimeout(5 * time.Second)},
	}
	client, err := zbc.NewClient(&config)
	if err != nil {
		return nil, fmt.Errorf("connect to camunda server error: %w", err)
	}

	return &CamundaClient{
		Client: client,
	}, nil
}

func (cli *CamundaClient) Close() error {
	return cli.Client.Close()
}

func (cli *CamundaClient) DeployProcess(ctx context.Context, name string, processDefinition []byte) (*pb.ProcessMetadata, error) {
	command := cli.Client.NewDeployResourceCommand().AddResource(processDefinition, name)

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

func (cli *CamundaClient) StartProcessInstance(ctx context.Context, processId string, version int32, vars map[string]any) (*pb.CreateProcessInstanceResponse, error) {
	var step3 commands.CreateInstanceCommandStep3
	step2 := cli.Client.NewCreateInstanceCommand().BPMNProcessId(processId)
	if version < 1 {
		step3 = step2.LatestVersion()
	} else {
		step3 = step2.Version(version)
	}
	command, err := step3.VariablesFromMap(vars)
	if err != nil {
		return nil, err
	}

	process, err := command.Send(ctx)
	if err != nil {
		return nil, err
	}

	return process, nil
}

type Option func(*startOptions)

type startOptions struct {
	concurrency   int
	maxJobActives int
	timeout       time.Duration
	pollInternal  time.Duration
}

func defaultStartOptions() startOptions {
	return startOptions{
		concurrency:   2,
		maxJobActives: 2,
		timeout:       time.Second * 10,
		pollInternal:  time.Second * 10,
	}
}

func WithConcurrency(concurrency int) Option {
	return func(so *startOptions) {
		so.concurrency = concurrency
	}
}

func WithMaxJobsActive(m int) Option {
	return func(so *startOptions) {
		so.maxJobActives = m
	}
}

func WithTimeout(t time.Duration) Option {
	return func(so *startOptions) {
		so.timeout = t
	}
}

func WithPoolInterval(t time.Duration) Option {
	return func(so *startOptions) {
		so.pollInternal = t
	}
}

func (cli *CamundaClient) StartWorker(jobType, workerName string, jobHandler worker.JobHandler, opts ...Option) worker.JobWorker {
	so := defaultStartOptions()
	for _, o := range opts {
		o(&so)
	}

	worker := cli.Client.NewJobWorker().JobType(jobType).Handler(jobHandler).Concurrency(so.concurrency).MaxJobsActive(so.maxJobActives).RequestTimeout(so.timeout).PollInterval(so.pollInternal).Name(workerName).Open()
	return worker
}
