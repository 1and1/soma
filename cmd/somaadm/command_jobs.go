package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/util"
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
								Usage:  `Check and update status of oustanding locally cached jobs`,
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
	if resp, err := adm.GetReq(`/jobs/`); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdJobShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if !utl.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s", c.Args().First())
	}
	path := fmt.Sprintf("/jobs/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	res, err := utl.ResultFromResponse(resp)
	if se, ok := err.(util.SomaError); ok {
		if se.RequestError() {
			return fmt.Errorf("Job update request error: %s", se.Error())
		}
		if se.Code() == 404 {
			return fmt.Errorf(`Could not find requested Job IDs`)
		}
		return fmt.Errorf("Job update application error: %s", err.Error())
	}
	for _, j := range *res.Jobs {
		if j.Status != `processed` {
			// only finish Jobs in DB that actually finished
			continue
		}
		strID := jobMap[j.Id]
		storeID, _ := strconv.ParseUint(strID, 10, 64)
		if err := store.FinishJob(storeID, &j); err != nil {
			return fmt.Errorf("somaadm: Job update cache error: %s", err.Error())
		}
	}
	fmt.Println(resp)
	return nil
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
	return nil
}

func cmdJobLocalPrune(c *cli.Context) error {
	return store.PruneFinishedJobs()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
