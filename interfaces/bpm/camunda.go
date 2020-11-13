package bpm

import (
	"errors"
	"fmt"
	camunda "github.com/citilinkru/camunda-client-go"
	"github.com/citilinkru/camunda-client-go/processor"
	"github.com/lanemets/claim-funnel/model"
	"log"
	"os"
)

type Camunda struct {
	client *camunda.Client
}

func NewCamundaClient(camunda *camunda.Client) Camunda {
	client := camunda
	return Camunda{
		client: client,
	}
}

func (c Camunda) RegisterExternalTaskWorker(
	workerConfig *WorkerConfig,
	handler ServiceTaskHandler,
) error {
	//TODO: multiple workers with different ids
	proc := createProcessor(c.client, workerConfig, handler.WorkerId)
	registerExternalHandler(proc, workerConfig, handler)

	return nil
}

func (c Camunda) DeployProcess(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error read file: %s\n", err)
		return errors.New("error read file")
	}
	result, err := c.client.Deployment.Create(
		camunda.ReqDeploymentCreate{
			DeploymentName: "DemoProcess",
			Resources: map[string]interface{}{
				"claim-process.bpmn": file,
			},
		})

	if err != nil {
		log.Printf("Error deploy process: %s\n", err)
		return errors.New("error deploy process")
	}

	log.Printf("process deployed successfully: %#+v\n", result)
	return nil
}

func (c Camunda) StartProcessInstance(processKey string, businessKey string) (*model.ProcessDefinitionId, error) {
	result, err := c.client.ProcessDefinition.StartInstance(
		camunda.QueryProcessDefinitionBy{Key: &processKey},
		camunda.ReqStartInstance{
			//TODO add variables
			//Variables: nil
			BusinessKey: &businessKey,
		},
	)
	if err != nil {
		msg := fmt.Sprintf("error start process: %s\n", err)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	log.Printf(
		"process instance has been started successfully; process definition id: %s, businessKey: %s",
		result.Id,
		result.BusinessKey,
	)

	return &model.ProcessDefinitionId{Value: result.DefinitionId}, nil
}

func (c Camunda) CompleteUserTask(businessKey string, taskId string, processDefinitionId *model.ProcessDefinitionId) error {
	//TODO: check for multiple process instances
	//TODO: retrieve task by query by id
	query := &camunda.UserTaskGetListQuery{
		ProcessInstanceBusinessKey: businessKey,
		ProcessDefinitionId:        processDefinitionId.Value,
		TaskDefinitionKey:          taskId,
	}

	tasks, err := c.client.UserTask.GetList(query)
	if err != nil {
		msg := fmt.Sprintf("unable to get user tasks; businessKey: %s, error: %s", businessKey, err)
		log.Println(msg)
		return errors.New(msg)
	}

	if len(tasks) == 0 {
		msg := fmt.Sprintf("no user task found; businessKey: %s, error: %s", businessKey, err)
		log.Println(msg)
		return errors.New(msg)
	}

	for _, t := range tasks {
		err := t.Complete(camunda.QueryUserTaskComplete{})
		//TODO: repeat? do not fail?
		if err != nil {
			//TODO: get rid of duplicate
			msg := fmt.Sprintf(
				"unable to complete user task; businessKey: %s, taskId: %s, error: %s",
				businessKey,
				taskId,
				err,
			)
			log.Println(msg)
		}
	}

	return nil
}

func createProcessor(client *camunda.Client, config *WorkerConfig, workerId string) *processor.Processor {
	return processor.NewProcessor(
		client,
		&processor.ProcessorOptions{
			WorkerId:                  workerId,
			LockDuration:              config.LockDuration,
			MaxTasks:                  config.MaxTasks,
			MaxParallelTaskPerHandler: config.MaxParallelTaskPerHandler,
			LongPollingTimeout:        config.LongPollingTimeout,
		},
		func(err error) {
			log.Println(err.Error())
		},
	)
}

func registerExternalHandler(
	p *processor.Processor,
	workerConfig *WorkerConfig,
	handler ServiceTaskHandler,
) {
	p.AddHandler(
		&[]camunda.QueryFetchAndLockTopic{
			{
				TopicName: handler.Topic,
			},
		},
		func(ctx *processor.Context) error {
			log.Printf(
				"Running task %s. WorkerId: %s. TopicName: %s. BusinessKey: %s \n",
				ctx.Task.Id,
				ctx.Task.WorkerId,
				ctx.Task.TopicName,
				ctx.Task.BusinessKey,
			)

			log.Printf("Topic: %v", ctx.Task.TopicName)
			log.Printf("Process Instance Id: %v", ctx.Task.ProcessInstanceId)
			log.Printf("Process Definition Id: %v", ctx.Task.ProcessDefinitionId)

			variablesToAdd, err := handler.Handler(variables(ctx.Task.Variables), ctx.Task.BusinessKey)

			log.Printf("Variables to add to the process: %v", variablesToAdd)

			if err != nil {
				errTxt := fmt.Sprintf(
					"an error occurred in executor; taskId: %s, businessKey: %s, error: %s",
					ctx.Task.Id,
					ctx.Task.BusinessKey,
					err,
				)

				log.Printf(errTxt)

				return ctx.HandleFailure(
					processor.QueryHandleFailure{
						ErrorMessage: &errTxt,
						Retries:      &workerConfig.Retries,
						RetryTimeout: &workerConfig.RetryTimeoutMillis,
					},
				)
			}

			log.Printf("Task variables: %v", ctx.Task.Variables)

			variables := camundaVariables(variablesToAdd)
			log.Printf("Variables to set: %v", variables)

			err = ctx.Complete(
				processor.QueryComplete{
					Variables: &variables,
				},
			)
			if err != nil {
				errTxt := fmt.Sprintf(
					"an error occurred on task completion; taskId: %s, businessKey: %s, error: %s",
					ctx.Task.Id,
					ctx.Task.BusinessKey,
					err,
				)

				log.Printf(errTxt)

				return ctx.HandleFailure(
					processor.QueryHandleFailure{
						ErrorMessage: &errTxt,
						Retries:      &workerConfig.Retries,
						RetryTimeout: &workerConfig.RetryTimeoutMillis,
					},
				)
			}

			log.Printf(
				"task completed successfully; taskId: %s, businessKey: %s",
				ctx.Task.Id,
				ctx.Task.BusinessKey,
			)
			return nil
		},
	)
}

func camundaVariables(attributes map[string]interface{}) map[string]camunda.Variable {
	instanceOf := func(i interface{}) string {
		switch t := i.(type) {
		case bool:
			return "boolean"
		case string:
			return "string"
		default:
			errorText := fmt.Sprintf("Don't know type %T\n", t)
			log.Fatalf(errorText)
			//TODO: remove
			return ""
		}
	}

	variables := make(map[string]camunda.Variable)

	for name, value := range attributes {
		kind := instanceOf(value)

		variables[name] = camunda.Variable{
			Value: value,
			Type:  kind,
		}
	}

	return variables
}

func variables(variables map[string]camunda.Variable) map[string]string {
	attributes := make(map[string]string)

	for name, value := range variables {
		attributes[name] = fmt.Sprintf("%v", value)
	}

	return attributes
}
