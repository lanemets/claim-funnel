package bpm

import (
	"errors"
	"fmt"
	camunda "github.com/citilinkru/camunda-client-go"
	"github.com/citilinkru/camunda-client-go/processor"
	"github.com/lanemets/claim-funnel/task"
	"log"
	"os"
	"time"
)

type Camunda struct {
	client *camunda.Client
}

type Credentials struct {
	endpointUrl string
	user        string
	password    string
}

type WorkerId struct {
	value string
}

type ProcessDefinitionId struct {
	value string
}

func NewCamundaClient(credentials Credentials) BpmClient {
	c := camunda.NewClient(
		camunda.ClientOptions{
			EndpointUrl: credentials.endpointUrl,
			ApiUser:     credentials.user,
			ApiPassword: credentials.password,
			Timeout:     time.Second * 10,
		},
	)

	return Camunda{
		client: c,
	}
}

//TODO: configure explicitly from external source
func configureExternalHandlers(p *processor.Processor, workerConfig *WorkerConfig, consumer task.Consumer) {
	p.AddHandler(
		&[]camunda.QueryFetchAndLockTopic{
			{
				TopicName: workerConfig.topicName,
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

			_ = consumer(ctx.Task.BusinessKey)

			err := ctx.Complete(
				processor.QueryComplete{
					Variables: &map[string]camunda.Variable{
						"result": {Value: "Hello world!", Type: "string"},
					},
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
						Retries:      &workerConfig.retries,
						RetryTimeout: &workerConfig.retryTimeout,
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

//TODO: move to separate goroutine/service
func (client Camunda) startExternalTaskWorker(config *WorkerConfig, consumer task.Consumer) error {

	//TODO: multiple workers with different ids
	p := createProcessor(client.client, "claim-process-worker-id")
	configureExternalHandlers(p, config, consumer)

	return nil
}

func (client Camunda) deployProcess(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error read file: %s\n", err)
		return errors.New("error read file")
	}
	result, err := client.client.Deployment.Create(
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

func (client Camunda) startProcessInstance(processKey string, businessKey string) (error, *ProcessDefinitionId) {
	result, err := client.client.ProcessDefinition.StartInstance(
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
		return errors.New(msg), nil
	}

	log.Printf(
		"process instance has been started successfully; process definition id: %s, businessKey: %s",
		result.Id,
		result.BusinessKey,
	)

	return nil, &ProcessDefinitionId{value: result.DefinitionId}
}

func (client Camunda) completeUserTask(businessKey string, processDefinitionId *ProcessDefinitionId) error {
	query := &camunda.UserTaskGetListQuery{
		ProcessInstanceBusinessKey: businessKey,
		ProcessDefinitionId:        processDefinitionId.value,
	}

	tasks, err := client.client.UserTask.GetList(query)
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
		if err != nil {
			//TODO: get rid of duplicate
			msg := fmt.Sprintf("unable to complete user task; businessKey: %s, error: %s", businessKey, err)
			log.Println(msg)
			return errors.New(msg)
		}
	}



	return nil
}

func createProcessor(client *camunda.Client, workerId string) *processor.Processor {
	//TODO: multiple workers?
	return processor.NewProcessor(
		client,
		&processor.ProcessorOptions{
			//TODO; to config
			WorkerId:                  workerId,
			LockDuration:              time.Second * 5,
			MaxTasks:                  10,
			MaxParallelTaskPerHandler: 100,
			LongPollingTimeout:        time.Minute,
		},
		func(err error) {
			log.Println(err.Error())
		},
	)
}
