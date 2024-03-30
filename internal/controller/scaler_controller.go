package controller

import (
	"context"
	"errors"
	"time"

	scalersv1beta1 "github.com/Malek-Zaag/ScalerOperator/api/v1beta1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = log.Log.WithName("OPERATOR LOGGER")

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.WithValues("Request Namespace", req.Namespace, "Request Name", req.Name)
	log.Info("Reconcile loop starting")
	scaler := &scalersv1beta1.Scaler{}
	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		return ctrl.Result{}, nil
	}
	// conenct to cluster and watch for metrics
	connection, err := cluster_connect()
	if err != nil {
		panic(err.Error())
	}
	podMetrics := get_metrics(connection)
	_ = podMetrics
	replicas := scaler.Spec.Replicas
	err1 := r.ScaleOnOverload(scaler, podMetrics, replicas, ctx)
	if err1 != nil {
		return ctrl.Result{}, err1
	}
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

func (r *ScalerReconciler) ScaleOnOverload(scaler *scalersv1beta1.Scaler, podMetrics []v1beta1.PodMetrics, replicas int32, ctx context.Context) error {
	for _, podMetric := range podMetrics {
		// check the pod name is the same as the deployment requested in the Scaler resource def
		if (podMetric.GetName()) == scaler.Spec.Deployments[0].Name && (podMetric.GetNamespace()) == scaler.Spec.Deployments[0].Namespace {
			if (podMetric.Containers[0].Usage.Cpu().MilliValue() > 50) || ((podMetric.Containers[0].Usage.Memory().Value() / (1024 * 1024)) > 200) {
				dep := &v1.Deployment{}
				err := r.Get(ctx, types.NamespacedName{
					Namespace: scaler.Spec.Deployments[0].Namespace,
					Name:      scaler.Spec.Deployments[0].Name,
				}, dep)
				if err != nil {
					return err
				}
				dep.Spec.Replicas = &replicas
				error := r.Update(ctx, dep)
				if error != nil {
					scaler.Status.Status = scalersv1beta1.FAILED
					return err
				}
				scaler.Status.Status = scalersv1beta1.SUCCESS
				error = r.Status().Update(ctx, scaler)
				if error != nil {
					return err
				}
			} else {
				logger.Info("all your pods are working fine, no overload is happening")
				dep := &v1.Deployment{}
				err := r.Get(ctx, types.NamespacedName{
					Namespace: scaler.Spec.Deployments[0].Namespace,
					Name:      scaler.Spec.Deployments[0].Name,
				}, dep)
				if err != nil {
					return err
				}
				scale_in := int32(1)
				dep.Spec.Replicas = &scale_in
				error := r.Update(ctx, dep)
				if error != nil {
					scaler.Status.Status = scalersv1beta1.FAILED
					return err
				}
				scaler.Status.Status = scalersv1beta1.SUCCESS
				error = r.Status().Update(ctx, scaler)
				if error != nil {
					return err
				}
			}
		} else {
			logger.Error(errors.New("your deployment cannot be found"), "your deployment cannot be found")
		}
	}
	return nil
}

func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalersv1beta1.Scaler{}).
		Complete(r)
}

// startTime := scaler.Spec.Start
// endTime := scaler.Spec.End
// currentHour := time.Now().UTC().Hour()
// fmt.Printf("The current hour is %d \n", currentHour)
// if currentHour >= startTime && currentHour <= endTime {
// 	fmt.Println("We are here")
// 	for _, deploy := range scaler.Spec.Deployments {
// 		dep := &v1.Deployment{}
// 		err := r.Get(ctx, types.NamespacedName{
// 			Namespace: deploy.Namespace,
// 			Name:      deploy.Name,
// 		}, dep)
// 		if err != nil {
// 			return ctrl.Result{}, err
// 		}
// 		if dep.Spec.Replicas != &replicas {
// 			dep.Spec.Replicas = &replicas
// 			err := r.Update(ctx, dep)
// 			if err != nil {
// 				scaler.Status.Status = scalersv1beta1.FAILED
// 				return ctrl.Result{}, err
// 			}
// 			scaler.Status.Status = scalersv1beta1.SUCCESS
// 			err = r.Status().Update(ctx, scaler)
// 			if err != nil {
// 				return ctrl.Result{}, err
// 			}
// 		}
// 	}
// }
