package main

type ErrorMarker interface {
	ErrorMark(err error, imp bool, found bool, length int, jobid string) bool
}

type SomaAppender interface {
	SomaAppendResult(r *somaResult)
	SomaAppendError(r *somaResult, err error)
}

type somaResult struct {
	RequestError   error
	NotFound       bool
	NotImplemented bool
	Accepted       bool
	JobId          string
	Nodes          []somaNodeResult
	Servers        []somaServerResult
	Levels         []somaLevelResult
	Predicates     []somaPredicateResult
	Status         []somaStatusResult
	Teams          []somaTeamResult
	Oncall         []somaOncallResult
	Views          []somaViewResult
	Units          []somaUnitResult
	Providers      []somaProviderResult
	Metrics        []somaMetricResult
	Modes          []somaModeResult
	Users          []somaUserResult
	Systems        []somaMonitoringResult
	Capabilities   []somaCapabilityResult
	Properties     []somaPropertyResult
	Attributes     []somaAttributeResult
	Repositories   []somaRepositoryResult
	Buckets        []somaBucketResult
	Groups         []somaGroupResult
	Clusters       []somaClusterResult
	CheckConfigs   []somaCheckConfigResult
	Deployments    []somaDeploymentResult
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
		ResultLength(r, reply), r.JobId)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
