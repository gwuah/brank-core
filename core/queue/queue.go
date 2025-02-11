package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"brank/core"
	"brank/core/storage"

	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

type JobIdentifier string

func (i JobIdentifier) String() string {
	return string(i)
}

type (
	JobWorker interface {
		Identifier() JobIdentifier
		Worker() que.WorkFunc
	}
	QueImpl interface {
		Close()
		RegisterJobs(jobList []JobWorker) *que.WorkerPool
		QueueJob(jobType JobIdentifier, payload interface{}) error
		QueueFutureJob(jobType JobIdentifier, payload interface{}, time ...time.Time) error
	}
	Que struct {
		dbURI    string
		config   *core.Config
		client   *que.Client
		connPool *pgx.ConnPool
	}
)

func getPgxPool(dbUri string) (*pgx.ConnPool, error) {
	pgxcfg, err := pgx.ParseURI(dbUri)
	if err != nil {
		return nil, err
	}
	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})
	if err != nil {
		return nil, err
	}
	return pgxpool, nil
}

func New(c *core.Config) (*Que, error) {
	q := &Que{dbURI: storage.GeneratePostgresURI(c), config: c}
	pgxpool, err := getPgxPool(q.dbURI)
	if err != nil {
		return nil, err
	}
	q.connPool = pgxpool
	q.client = que.NewClient(pgxpool)
	return q, nil
}

func (q *Que) Close() {
	log.Println("shutting down queue")
	q.connPool.Close()
}

func (q *Que) RegisterJobs(jobList []JobWorker) *que.WorkerPool {
	wm := que.WorkMap{}
	for _, j := range jobList {
		wm[j.Identifier().String()] = j.Worker()
	}
	return que.NewWorkerPool(q.client, wm, q.config.WORKER_POOL_SIZE)
}

func (q *Que) QueueJob(jobType JobIdentifier, payload interface{}) error {
	enc, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	j := que.Job{Type: jobType.String(), Args: enc}
	err = q.client.Enqueue(&j)
	if err != nil {
		return fmt.Errorf("failed to queue job. err: %w", err)
	}
	return nil
}

func (q *Que) QueueFutureJob(jobType JobIdentifier, payload interface{}, times ...time.Time) error {
	enc, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for _, time := range times {
		j := que.Job{Type: jobType.String(), Args: enc, RunAt: time}
		err = q.client.Enqueue(&j)
		if err != nil {
			return fmt.Errorf("failed to queue job. err: %w", err)
		}
	}

	return nil
}
