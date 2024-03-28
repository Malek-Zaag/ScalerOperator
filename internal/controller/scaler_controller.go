package controller

import (
	"context"
	"fmt"
	"time"

	scalersv1beta1 "github.com/Malek-Zaag/ScalerOperator/api/v1beta1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/metrics"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = log.Log.WithName("controller_scaler")

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	log.Info("Reconcile loop")
	// TODO(user): your logic here
	scaler := &scalersv1beta1.Scaler{}
	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		return ctrl.Result{}, nil
	}
	connection, err := cluster_connect()
	if err != nil {
		panic(err.Error())
	}
	get_metrics(connection)
	startTime := scaler.Spec.Start
	endTime := scaler.Spec.End
	replicas := scaler.Spec.Replicas
	fmt.Print(metrics.GroupName)
	_ = replicas
	currentHour := time.Now().UTC().Hour()
	fmt.Printf("The current hour is %d \n", currentHour)
	if currentHour >= startTime && currentHour <= endTime {
		fmt.Println("We are here")
		for _, deploy := range scaler.Spec.Deployments {
			dep := &v1.Deployment{}
			err := r.Get(ctx, types.NamespacedName{
				Namespace: deploy.Namespace,
				Name:      deploy.Name,
			}, dep)
			if err != nil {
				return ctrl.Result{}, err
			}
			if dep.Spec.Replicas != &replicas {
				dep.Spec.Replicas = &replicas
				err := r.Update(ctx, dep)
				if err != nil {
					scaler.Status.Status = scalersv1beta1.FAILED
					return ctrl.Result{}, err
				}
				scaler.Status.Status = scalersv1beta1.SUCCESS
				err = r.Status().Update(ctx, scaler)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalersv1beta1.Scaler{}).
		Complete(r)
}
