/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and dns-masquerading-operator-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package operator

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"time"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	istionetworkingv1 "istio.io/client-go/pkg/apis/networking/v1"

	"github.com/sap/component-operator-runtime/pkg/component"
	helmgenerator "github.com/sap/component-operator-runtime/pkg/manifests/helm"
	"github.com/sap/component-operator-runtime/pkg/operator"
	componentoperatorruntimetypes "github.com/sap/component-operator-runtime/pkg/types"

	operatorv1alpha1 "github.com/sap/dns-masquerading-operator-cop/api/v1alpha1"
	"github.com/sap/dns-masquerading-operator-cop/internal/transformer"
)

const Name = "dns-masquerading-operator-cop.cs.sap.com"

const (
	managedIndexKey = ".metadata.managed"
	finalizer       = "dns.cs.sap.com/masquerading-operator"
)

//go:embed all:data
var data embed.FS

type Options struct {
	Name                  string
	DefaultServiceAccount string
	FlagPrefix            string
}

type Operator struct {
	options Options
}

var defaultOperator operator.Operator = New()

func GetName() string {
	return defaultOperator.GetName()
}

func InitScheme(scheme *runtime.Scheme) {
	defaultOperator.InitScheme(scheme)
}

func InitFlags(flagset *flag.FlagSet) {
	defaultOperator.InitFlags(flagset)
}

func ValidateFlags() error {
	return defaultOperator.ValidateFlags()
}

func GetUncacheableTypes() []client.Object {
	return defaultOperator.GetUncacheableTypes()
}

func Setup(mgr ctrl.Manager) error {
	return defaultOperator.Setup(mgr)
}

func New() *Operator {
	return NewWithOptions(Options{})
}

func NewWithOptions(options Options) *Operator {
	operator := &Operator{options: options}
	if operator.options.Name == "" {
		operator.options.Name = Name
	}
	return operator
}

func (o *Operator) GetName() string {
	return o.options.Name
}

func (o *Operator) InitScheme(scheme *runtime.Scheme) {
	utilruntime.Must(istionetworkingv1.AddToScheme(scheme))
	utilruntime.Must(operatorv1alpha1.AddToScheme(scheme))
}

func (o *Operator) InitFlags(flagset *flag.FlagSet) {
	flagset.StringVar(&o.options.DefaultServiceAccount, "default-service-account", o.options.DefaultServiceAccount, "Default service account name")
}

func (o *Operator) ValidateFlags() error {
	return nil
}

func (o *Operator) GetUncacheableTypes() []client.Object {
	return []client.Object{&operatorv1alpha1.DNSMasqueradingOperator{}}
}

func (o *Operator) Setup(mgr ctrl.Manager) error {
	if err := mgr.GetCache().IndexField(context.TODO(), &corev1.Service{}, managedIndexKey, indexByManaged); err != nil {
		return errors.Wrapf(err, "failed setting index field %s", managedIndexKey)
	}
	if err := mgr.GetCache().IndexField(context.TODO(), &networkingv1.Ingress{}, managedIndexKey, indexByManaged); err != nil {
		return errors.Wrapf(err, "failed setting index field %s", managedIndexKey)
	}
	/*
		if err := mgr.GetCache().IndexField(context.TODO(), &istionetworkingv1.Gateway{}, managedIndexKey, indexByManaged); err != nil {
			return errors.Wrapf(err, "failed setting index field %s", managedIndexKey)
		}
	*/

	resourceGenerator, err := helmgenerator.NewHelmGeneratorWithParameterTransformer(
		data,
		"data/charts/dns-masquerading-operator",
		mgr.GetClient(),
		transformer.NewParameterTransformer(),
	)
	if err != nil {
		return errors.Wrap(err, "error initializing resource generator")
	}

	if err := component.NewReconciler[*operatorv1alpha1.DNSMasqueradingOperator](
		o.options.Name,
		resourceGenerator,
		component.ReconcilerOptions{
			DefaultServiceAccount: &o.options.DefaultServiceAccount,
		},
	).WithPreDeleteHook(preDeleteHook).SetupWithManager(mgr); err != nil {
		return errors.Wrapf(err, "unable to create controller")
	}

	return nil
}

func indexByManaged(object client.Object) []string {
	if controllerutil.ContainsFinalizer(object, finalizer) {
		return []string{"true"}
	}
	return nil
}

func preDeleteHook(ctx context.Context, clnt client.Client, c *operatorv1alpha1.DNSMasqueradingOperator) error {
	// TODO: it would be more elegant to maintain a custom cache index and use MatchingFields in the List() call ...
	// TODO: limit results to 1

	if c.Spec.EnableServiceController {
		serviceList := &corev1.ServiceList{}
		if err := clnt.List(ctx, serviceList, client.MatchingFields{managedIndexKey: "true"}); err != nil {
			return err
		}
		if len(serviceList.Items) > 0 {
			service := serviceList.Items[0]
			return componentoperatorruntimetypes.NewRetriableError(fmt.Errorf("deletion blocked by service %s/%s", service.Namespace, service.Name), &[]time.Duration{10 * time.Second}[0])
		}
	}

	if c.Spec.EnableIngressController {
		ingressList := &networkingv1.IngressList{}
		if err := clnt.List(ctx, ingressList, client.MatchingFields{managedIndexKey: "true"}); err != nil {
			return err
		}
		if len(ingressList.Items) > 0 {
			ingress := ingressList.Items[0]
			return componentoperatorruntimetypes.NewRetriableError(fmt.Errorf("deletion blocked by ingress %s/%s", ingress.Namespace, ingress.Name), &[]time.Duration{10 * time.Second}[0])
		}
	}

	if c.Spec.EnableIstioGatewayController {
		gatewayList := &istionetworkingv1.GatewayList{}
		if err := clnt.List(ctx, gatewayList, client.MatchingFields{managedIndexKey: "true"}); err != nil {
			return err
		}
		if len(gatewayList.Items) > 0 {
			gateway := gatewayList.Items[0]
			return componentoperatorruntimetypes.NewRetriableError(fmt.Errorf("deletion blocked by gateway %s/%s", gateway.Namespace, gateway.Name), &[]time.Duration{10 * time.Second}[0])
		}
	}

	return nil
}
