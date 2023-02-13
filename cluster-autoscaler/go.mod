module k8s.io/autoscaler/cluster-autoscaler

go 1.19

require (
	cloud.google.com/go/compute v1.7.0
	github.com/Azure/azure-sdk-for-go v65.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.24
	github.com/Azure/go-autorest/autorest/adal v0.9.18
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/aws/aws-sdk-go v1.38.49
	github.com/digitalocean/godo v1.27.0
	github.com/ghodss/yaml v1.0.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/google/go-querystring v1.1.0
	github.com/google/uuid v1.3.0
	github.com/jmespath/go-jmespath v0.4.0
	github.com/json-iterator/go v1.1.12
	github.com/juju/juju v0.0.0-20230209091443-f2792de1004b
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	golang.org/x/crypto v0.5.0
	golang.org/x/oauth2 v0.0.0-20221006150949-b44042a4b9c1
	google.golang.org/api v0.84.0
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/gcfg.v1 v1.2.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.24.0-alpha.4
	k8s.io/apimachinery v0.24.0-alpha.4
	k8s.io/apiserver v0.24.0-alpha.4
	k8s.io/client-go v0.24.0-alpha.4
	k8s.io/cloud-provider v0.24.0-alpha.4
	k8s.io/component-base v0.24.0-alpha.4
	k8s.io/component-helpers v0.24.0-alpha.4
	k8s.io/klog/v2 v2.40.1
	k8s.io/kubelet v0.23.0
	k8s.io/kubernetes v1.24.0-alpha.4
	k8s.io/legacy-cloud-providers v0.0.0
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/cloud-provider-azure v1.23.2
)

require (
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/mocks v0.4.1 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/GoogleCloudPlatform/k8s-cloud-provider v1.16.1-0.20210702024009-ea6160c1d0e3 // indirect
	github.com/JeffAshton/win_pdh v0.0.0-20161109143554-76bb4ee9f0ab // indirect
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/Microsoft/hcsshim v0.8.22 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e // indirect
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/aws/aws-sdk-go-v2 v1.9.1 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.2.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.2.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.4.1 // indirect
	github.com/aws/smithy-go v1.8.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/checkpoint-restore/go-criu/v5 v5.0.0 // indirect
	github.com/cilium/ebpf v0.6.2 // indirect
	github.com/clusterhq/flocker-go v0.0.0-20160920122132-2b8b7259d313 // indirect
	github.com/container-storage-interface/spec v1.5.0 // indirect
	github.com/containerd/cgroups v1.0.1 // indirect
	github.com/containerd/console v1.0.2 // indirect
	github.com/containerd/containerd v1.4.11 // indirect
	github.com/containerd/ttrpc v1.0.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/euank/go-kmsg-parser v2.0.0+incompatible // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-macaroon-bakery/macaroon-bakery/v3 v3.0.1 // indirect
	github.com/go-macaroon-bakery/macaroonpb v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godbus/dbus/v5 v5.0.4 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/cadvisor v0.43.0 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.0.0-20220520183353-fd19c99a87aa // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/heketi/heketi v10.3.0+incompatible // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/juju/ansiterm v1.0.0 // indirect
	github.com/juju/charm/v8 v8.0.6 // indirect
	github.com/juju/charmrepo/v6 v6.0.3 // indirect
	github.com/juju/clock v1.0.3 // indirect
	github.com/juju/cmd/v3 v3.0.4 // indirect
	github.com/juju/collections v1.0.2 // indirect
	github.com/juju/errors v1.0.0 // indirect
	github.com/juju/featureflag v1.0.0 // indirect
	github.com/juju/gnuflag v1.0.0 // indirect
	github.com/juju/go4 v0.0.0-20160222163258-40d72ab9641a // indirect
	github.com/juju/gojsonpointer v0.0.0-20150204194629-afe8b77aa08f // indirect
	github.com/juju/gojsonreference v0.0.0-20150204194633-f0d24ac5ee33 // indirect
	github.com/juju/gojsonschema v1.0.0 // indirect
	github.com/juju/http/v2 v2.0.0 // indirect
	github.com/juju/idmclient/v2 v2.0.0 // indirect
	github.com/juju/loggo v1.0.0 // indirect
	github.com/juju/mgo/v2 v2.0.2 // indirect
	github.com/juju/mgo/v3 v3.0.3 // indirect
	github.com/juju/mutex/v2 v2.0.0 // indirect
	github.com/juju/names/v4 v4.0.0 // indirect
	github.com/juju/os/v2 v2.2.3 // indirect
	github.com/juju/persistent-cookiejar v1.0.0 // indirect
	github.com/juju/proxy v1.0.0 // indirect
	github.com/juju/retry v1.0.0 // indirect
	github.com/juju/rfc/v2 v2.0.0 // indirect
	github.com/juju/romulus v1.0.0 // indirect
	github.com/juju/rpcreflect v1.0.0 // indirect
	github.com/juju/schema v1.0.1 // indirect
	github.com/juju/usso v1.0.1 // indirect
	github.com/juju/utils/v3 v3.0.2 // indirect
	github.com/juju/version/v2 v2.0.1 // indirect
	github.com/juju/webbrowser v1.0.0 // indirect
	github.com/juju/worker/v3 v3.1.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/libopenstorage/openstorage v1.0.0 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/lxc/lxd v0.0.0-20220816180258-7e0418163fa9 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989 // indirect
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/ipvs v1.0.1 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runc v1.0.3 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/opencontainers/selinux v1.8.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/quobyte/api v0.1.8 // indirect
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/rubiojr/go-vhd v0.0.0-20200706105327-02e210299021 // indirect
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/cobra v1.5.0 // indirect
	github.com/storageos/go-api v2.2.0+incompatible // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2 // indirect
	github.com/vishvananda/netns v0.0.0-20211101163701-50045581ed74 // indirect
	github.com/vmware/govmomi v0.21.1-0.20191008161538-40aebf13ba45 // indirect
	go.etcd.io/etcd/api/v3 v3.5.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib v0.20.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.20.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0 // indirect
	go.opentelemetry.io/otel v0.20.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp v0.20.0 // indirect
	go.opentelemetry.io/otel/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/trace v0.20.0 // indirect
	go.opentelemetry.io/proto/otlp v0.7.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sync v0.0.0-20220929204114-8fcdb60fdcc0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/term v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220720214146-176da50484ac // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/gobwas/glob.v0 v0.2.3 // indirect
	gopkg.in/httprequest.v1 v1.2.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/juju/environschema.v1 v1.0.1-0.20201027142642-c89a4490670a // indirect
	gopkg.in/macaroon.v2 v2.1.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/retry.v1 v1.0.3 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	gopkg.in/warnings.v0 v0.1.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/cri-api v0.0.0 // indirect
	k8s.io/csi-translation-lib v0.24.0-alpha.4 // indirect
	k8s.io/kube-openapi v0.0.0-20220316025549-ddc66922ab18 // indirect
	k8s.io/kube-proxy v0.0.0 // indirect
	k8s.io/kube-scheduler v0.0.0 // indirect
	k8s.io/kubectl v0.0.0 // indirect
	k8s.io/mount-utils v0.24.0-alpha.4 // indirect
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.0.30 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/aws/aws-sdk-go/service/eks => github.com/aws/aws-sdk-go/service/eks v1.38.49

replace github.com/digitalocean/godo => github.com/digitalocean/godo v1.27.0

replace github.com/rancher/go-rancher => github.com/rancher/go-rancher v0.1.0

replace k8s.io/api => k8s.io/api v0.24.0-alpha.4

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.24.0-alpha.4

replace k8s.io/apimachinery => k8s.io/apimachinery v0.24.0-alpha.4

replace k8s.io/apiserver => k8s.io/apiserver v0.24.0-alpha.4

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.24.0-alpha.4

replace k8s.io/client-go => k8s.io/client-go v0.24.0-alpha.4

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.24.0-alpha.4

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.24.0-alpha.4

replace k8s.io/code-generator => k8s.io/code-generator v0.24.0-alpha.4

replace k8s.io/component-base => k8s.io/component-base v0.24.0-alpha.4

replace k8s.io/component-helpers => k8s.io/component-helpers v0.24.0-alpha.4

replace k8s.io/controller-manager => k8s.io/controller-manager v0.24.0-alpha.4

replace k8s.io/cri-api => k8s.io/cri-api v0.24.0-alpha.4

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.24.0-alpha.4

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.24.0-alpha.4

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.24.0-alpha.4

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.24.0-alpha.4

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.24.0-alpha.4

replace k8s.io/kubectl => k8s.io/kubectl v0.24.0-alpha.4

replace k8s.io/kubelet => k8s.io/kubelet v0.24.0-alpha.4

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.24.0-alpha.4

replace k8s.io/metrics => k8s.io/metrics v0.24.0-alpha.4

replace k8s.io/mount-utils => k8s.io/mount-utils v0.24.0-alpha.4

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.24.0-alpha.4

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.24.0-alpha.4

replace k8s.io/sample-controller => k8s.io/sample-controller v0.24.0-alpha.4

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.24.0-alpha.4
