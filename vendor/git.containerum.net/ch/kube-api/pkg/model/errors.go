package model

import (
	"errors"

	ch "git.containerum.net/ch/kube-client/pkg/cherry"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/kube-api"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	api_meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrInvalidCPUFormat    = errors.New("invalid cpu quota format")
	ErrInvalidMemoryFormat = errors.New("invalid memory quota format")

	ErrUnableDecodeUserHeaderData    = errors.New("unable to decode user header data")
	ErrUnableUnmarshalUserHeaderData = errors.New("unable to unmarshal user header data")

	ErrUnableConvertServiceList = errors.New("unable to decode services list")
	ErrUnableConvertService     = errors.New("unable to decode service")

	ErrUnableConvertNamespaceList = errors.New("unable to decode namespaces list")
	ErrUnableConvertNamespace     = errors.New("unable to decode namespace")

	ErrUnableConvertSecretList = errors.New("unable to decode secrets list")
	ErrUnableConvertSecret     = errors.New("unable to decode secret")

	ErrUnableConvertIngressList = errors.New("unable to decode ingresses list")
	ErrUnableConvertIngress     = errors.New("unable to decode ingress")

	ErrUnableConvertDeploymentList = errors.New("unable to decode deployment list")
	ErrUnableConvertDeployment     = errors.New("unable to decode deployment")

	ErrUnableConvertEndpointList = errors.New("unable to decode services list")
	ErrUnableConvertEndpoint     = errors.New("unable to decode service")

	ErrUnableConvertConfigMapList = errors.New("unable to decode config maps list")
	ErrUnableConvertConfigMap     = errors.New("unable to decode config map")
)

const (
	noContainer         = "container %v is not found in deployment"
	fieldShouldExist    = "field %v should be provided"
	invalidReplicas     = "invalid replicas number: %v. It must be between 1 and %v"
	invalidPort         = "invalid port: %v. It must be between %v and %v"
	invalidProtocol     = "invalid protocol: %v. It must be TCP or UDP"
	invalidOwner        = "owner should be UUID"
	invalidName         = "invalid name: %v. %v"
	invalidIP           = "invalid IP: %v. It must be a valid IP address, (e.g. 10.9.8.7)"
	invalidCPUQuota     = "invalid CPU quota: %v. It must be between %v and %v"
	invalidMemoryQuota  = "invalid memory quota: %v. It must be between %v and %v"
	subPathRelative     = "invalid Sub Path: %v. It must be relative path"
	invalidResourceKind = "invalid resource kind: %v. Shoud be %v"
	invalidAPIVersion   = "invalid API Version: %v. Shoud be %v"
)

//ParseResourceError checks error status
func ParseResourceError(in interface{}, defaulterr *ch.Err) *ch.Err {
	sE, isStatusErrorCode := in.(*api_errors.StatusError)
	if isStatusErrorCode {
		switch sE.ErrStatus.Reason {
		case api_meta.StatusReasonNotFound:
			return cherry.ErrResourceNotExist()
		case api_meta.StatusReasonAlreadyExists:
			return cherry.ErrResourceAlreadyExists()
		default:
			return defaulterr
		}
	}
	return defaulterr
}
