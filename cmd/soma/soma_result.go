package main

type ErrorMarker interface {
	ErrorMark(err error, imp bool, found bool, length int, jobid, jobtype string) bool
}

type SomaAppender interface {
	SomaAppendResult(r *somaResult)
	SomaAppendError(r *somaResult, err error)
}

type somaResult struct {
	RequestError    error
	NotFound        bool
	NotImplemented  bool
	Accepted        bool
	JobId           string
	JobType         string
	Attributes      []somaAttributeResult
	Buckets         []somaBucketResult
	Capabilities    []somaCapabilityResult
	CheckConfigs    []somaCheckConfigResult
	Clusters        []somaClusterResult
	Datacenters     []somaDatacenterResult
	Deployments     []somaDeploymentResult
	Groups          []somaGroupResult
	HostDeployments []somaHostDeploymentResult
	Levels          []somaLevelResult
	Metrics         []somaMetricResult
	Modes           []somaModeResult
	Nodes           []somaNodeResult
	Oncall          []somaOncallResult
	Predicates      []somaPredicateResult
	Properties      []somaPropertyResult
	Providers       []somaProviderResult
	Repositories    []somaRepositoryResult
	Servers         []somaServerResult
	Status          []somaStatusResult
	Systems         []somaMonitoringResult
	Teams           []somaTeamResult
	Units           []somaUnitResult
	Users           []somaUserResult
	Validity        []somaValidityResult
	Views           []somaViewResult
}

func (r *somaResult) SetRequestError(err error) bool {
	if err != nil {
		r.RequestError = err
		return true
	}
	return false
}

func (r *somaResult) SetNotFound() {
	r.NotFound = true
}

func (r *somaResult) SetNotFoundErr(err error) {
	r.NotFound = true
	r.RequestError = err
}

func (r *somaResult) SetNotImplemented() {
	r.NotImplemented = true
}

func (r *somaResult) Failure() bool {
	if r.NotFound || r.NotImplemented || r.RequestError != nil {
		return true
	}
	return false
}

func (r *somaResult) Append(err error, res SomaAppender) {
	if err != nil {
		res.SomaAppendError(r, err)
		return
	}
	res.SomaAppendResult(r)
}

func (r *somaResult) MarkErrors(reply ErrorMarker) bool {
	return reply.ErrorMark(r.RequestError, r.NotImplemented, r.NotFound,
		ResultLength(r, reply), r.JobId, r.JobType)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
