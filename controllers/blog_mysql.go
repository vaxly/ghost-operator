package controllers

import (
	"context"

	ghostv1alpha1 "github.com/vaxly/ghost-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *BlogReconciler) MysqlCreateOrUpdatePersistentVolumeClaim(cr *ghostv1alpha1.Blog) error {
	requestStorage := make(corev1.ResourceList)
	requestStorage[corev1.ResourceStorage] = cr.Spec.Persistent.Size
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlPersistentVolumeClaimNameFromCR(cr),
			Namespace: cr.Namespace,
			Labels:    mysqlLabelFromCR(cr),
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, pvc, func() error {
		if err := controllerutil.SetControllerReference(cr, pvc, r.Scheme); err != nil {
			return err
		}

		if pvc.ObjectMeta.CreationTimestamp.IsZero() {
			pvc.Spec = corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				StorageClassName: cr.Spec.Persistent.StorageClass,
			}
		}

		pvc.Spec.Resources = corev1.ResourceRequirements{
			Requests: requestStorage,
		}

		return nil
	})

	r.Logger.Info("Reconciling MySql PersistentVolumeClaim ", pvc.Name)
	return err
}

func (r *BlogReconciler) MyqlsCreateOrUpdateService(cr *ghostv1alpha1.Blog) error {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + mysqlSuffix,
			Namespace: cr.Namespace,
			Labels:    mysqlLabelFromCR(cr),
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, svc, func() error {
		// We don't accept any update for service
		if !svc.ObjectMeta.CreationTimestamp.IsZero() {
			return nil
		}

		if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
			return err
		}

		svc.Spec = corev1.ServiceSpec{
			Selector: mysqlLabelFromCR(cr),
			Type:     cr.Spec.ServiceType,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   "TCP",
					Port:       int32(mysqlPort),
					TargetPort: intstr.FromInt(mysqlPort),
				},
			},
		}
		return nil
	})

	r.Logger.Info("Reconciling MySql Service ", svc.Name)
	return err
}

func (r *BlogReconciler) MysqlCreateOrUpdateDeployment(cr *ghostv1alpha1.Blog) error {
	defaultTerminationGracePeriodSeconds := int64(30)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + mysqlSuffix,
			Namespace: cr.Namespace,
			Labels:    mysqlLabelFromCR(cr),
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, dep, func() error {
		if dep.ObjectMeta.CreationTimestamp.IsZero() {
			// Set label selector only when deployment has never been created
			dep.Spec.Selector = mysqlLabelSelectorFromCR(cr)
		}

		if err := controllerutil.SetControllerReference(cr, dep, r.Scheme); err != nil {
			return err
		}

		dep.Spec.Replicas = cr.Spec.Replicas
		dep.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: mysqlLabelFromCR(cr),
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "mysql",
						Image:           "mysql:5.7",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Env: []corev1.EnvVar{{
							Name:  "MYSQL_ROOT_PASSWORD",
							Value: cr.Spec.Config.Database.Connection.Password,
						}},
						Ports: []corev1.ContainerPort{
							{
								Name:          "mysql",
								ContainerPort: int32(mysqlPort),
								Protocol:      corev1.ProtocolTCP,
							},
						},
						TerminationMessagePath:   "/dev/termination-log",
						TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mysql-persistent-storage",
							MountPath: "/var/lib/mysql",
						}},
					},
				},
				RestartPolicy:                 corev1.RestartPolicyAlways,
				TerminationGracePeriodSeconds: &defaultTerminationGracePeriodSeconds,
				DNSPolicy:                     corev1.DNSClusterFirst,
				SecurityContext:               &corev1.PodSecurityContext{},
				SchedulerName:                 corev1.DefaultSchedulerName,
				Volumes: []corev1.Volume{{
					Name: "mysql-persistent-storage",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: mysqlPersistentVolumeClaimNameFromCR(cr),
						},
					},
				}},
			},
		}
		return nil
	})

	r.Logger.Info("Reconciling MySql Deployment ", dep.Name)
	return err
}

func (r *BlogReconciler) ISMysqlUp(cr *ghostv1alpha1.Blog) bool {
	deployment := &appsv1.Deployment{}

	if err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cr.Name + mysqlSuffix,
		Namespace: cr.Namespace,
	}, deployment); err != nil {
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false

}
