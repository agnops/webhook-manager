module github.com/agnops/webhook-manager/app

go 1.15

require (
	github.com/google/go-github v17.0.0+incompatible
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/xanzy/go-gitlab v0.32.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/appengine v1.6.6 // indirect
	k8s.io/api v0.0.0-20200726131424-9540e4cac147
	k8s.io/apimachinery v0.0.0-20200726131235-945d4ebf362b
	k8s.io/client-go v0.0.0-20200726131703-36233866f1c7
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.2.0
	k8s.io/utils v0.0.0-20200720150651-0bdb4ca86cbc
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20200726131424-9540e4cac147
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200726131235-945d4ebf362b
)
