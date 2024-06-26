package builder

import (
	mariadbv1alpha1 "github.com/mariadb-operator/mariadb-operator/api/v1alpha1"
	"github.com/mariadb-operator/mariadb-operator/pkg/environment"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func newTestBuilder() *Builder {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mariadbv1alpha1.AddToScheme(scheme))
	utilruntime.Must(monitoringv1.AddToScheme(scheme))

	return NewBuilder(scheme, &environment.OperatorEnv{
		MariadbOperatorName:      "mariadb-operator",
		MariadbOperatorNamespace: "test",
		MariadbOperatorSAPath:    "/var/run/secrets/kubernetes.io/serviceaccount/token",
		MariadbOperatorImage:     "mariadb-operator:test",
		RelatedMariadbImage:      "mariadb:11.2.2:test",
		RelatedMaxscaleImage:     "maxscale:test",
		RelatedExporterImage:     "mysql-exporter:test",
		MariadbGaleraInitImage:   "mariadb-operator:tes",
		MariadbGaleraAgentImage:  "mariadb-operator:test",
		MariadbGaleraLibPath:     "/usr/lib/galera/libgalera_smm.so",
		WatchNamespace:           "",
	})
}
