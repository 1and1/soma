package main

import (
	"encoding/json"
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

func registerJobs(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `jobs`,
				Usage: `SUBCOMMANDS for job information`,
				Subcommands: []cli.Command{
					{
						Name:   `list`,
						Usage:  `List outstanding jobs (remote)`,
						Action: runtime(cmdJobList),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "all, a",
								Usage: "List all outstanding jobs (admin only)",
							},
						},
					},
					{
						Name:   `show`,
						Usage:  `Show details about a job (remote)`,
						Action: runtime(cmdJobShow),
					},
					{
						Name:  `local`,
						Usage: `SUBCOMMANDS for locally saved jobs`,
						Subcommands: []cli.Command{
							{
								Name:   `outstanding`,
								Usage:  `List outstanding locally saved Jobs`,
								Action: runtime(cmdJobLocalOutstanding),
							},
							{
								Name:   `update`,
								Usage:  `Check and update status of outstanding locally cached jobs`,
								Action: runtime(cmdJobLocalUpdate),
								Flags: []cli.Flag{
									cli.BoolFlag{
										Name:  "verbose, v",
										Usage: "Include full raw job request (admin only)",
									},
								},
							},
							{
								Name:   `list`,
								Usage:  `List all locally cached jobs`,
								Action: runtime(cmdJobLocalList),
							},
							{
								Name:   `prune`,
								Usage:  `Delete completed jobs from local cache`,
								Action: runtime(cmdJobLocalPrune),
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

func cmdJobList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if c.Bool(`all`) {
		return adm.Perform(`get`, `/jobs/all`, `list`, nil, c)
	}
	return adm.Perform(`get`, `/jobs/`, `list`, nil, c)
}

func cmdJobShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/jobs/id/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdJobLocalOutstanding(c *cli.Context) error {
	jobs, err := store.ActiveJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	pj := []proto.Job{}
	for _, iArray := range jobs {
		pj = append(pj, proto.Job{
			Id:       iArray[1],
			TsQueued: iArray[2],
			Type:     iArray[3],
		})
	}

	if enc, err := json.Marshal(&pj); err != nil {
		return err
	} else {
		fmt.Println(string(enc))
	}
	// XXX adm.FormatOut support missing
	return nil
}

func cmdJobLocalUpdate(c *cli.Context) error {
	jobs, err := store.ActiveJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	} else if err == bolt.ErrBucketNotFound {
		// nothing found
		return nil
	}

	req := proto.NewJobFilter()
	req.Flags.Detailed = c.Bool(`verbose`)
	jobMap := map[string]string{}
	for _, v := range jobs {
		// jobID -> storeID
		jobMap[v[1]] = v[0]
		req.Filter.Job.IdList = append(req.Filter.Job.IdList, v[1])
	}
	resp, err := adm.PostReqBody(req, `/filter/jobs/`)
	if err != nil {
		return fmt.Errorf("Job update request error: %s", err)
	}
	var res *proto.Result
	if err = adm.DecodedResponse(resp, res); err != nil {
		return err
	}
	if res.Jobs == nil {
		return fmt.Errorf("Result contained no jobs array")
	}
	for _, j := range *res.Jobs {
		if j.Status != `processed` {
			// only finish Jobs in DB that actually finished
			continue
		}
		strID := jobMap[j.Id]
		var storeID uint64
		if err := adm.ValidateLBoundUint64(strID, &storeID,
			0); err != nil {
			return fmt.Errorf("somaadm: Job update cache error: %s",
				err.Error())
		}
		if err := store.FinishJob(storeID, &j); err != nil {
			return fmt.Errorf("somaadm: Job update cache error: %s",
				err.Error())
		}
	}
	return adm.FormatOut(c, resp, `list`)
}

func cmdJobLocalList(c *cli.Context) error {
	active, err := store.ActiveJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	jobs := []proto.Job{}
	for _, iArray := range active {
		jobs = append(jobs, proto.Job{
			Id:       iArray[1],
			TsQueued: iArray[2],
			Type:     iArray[3],
		})
	}

	finished, err := store.FinishedJobs()
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	jobs = append(jobs, finished...)
	if enc, err := json.Marshal(&jobs); err != nil {
		return err
	} else {
		fmt.Println(string(enc))
	}
	// XXX adm.FormatOut support missing
	return nil
}

func cmdJobLocalPrune(c *cli.Context) error {
	return store.PruneFinishedJobs()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
