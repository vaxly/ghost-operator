package controllers

import (
	"context"
	"encoding/json"

	ghostv1alpha1 "github.com/vaxly/ghost-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *BlogReconciler) newVolumeForCR(cr *ghostv1alpha1.Blog) []corev1.Volume {
	configMapDefaultMode := int32(0644)
	var volume []corev1.Volume
	volume = append(volume, corev1.Volume{
		Name: "ghost-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapNameFromCR(cr),
				},
				DefaultMode: &configMapDefaultMode,
			},
		},
	})

	var ghostContentSource corev1.VolumeSource
	if cr.Spec.Persistent.Enabled {
		ghostContentSource = corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: persistentVolumeClaimNameFromCR(cr),
			},
		}
	} else {
		ghostContentSource = corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		}
	}

	volume = append(volume, corev1.Volume{
		Name:         "ghost-content",
		VolumeSource: ghostContentSource,
	})

	return volume
}

func (r *BlogReconciler) CreateOrUpdateIngress(cr *ghostv1alpha1.Blog) error {
	ing := &networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.GetName(),
			Namespace:   cr.GetNamespace(),
			Labels:      commonLabelFromCR(cr),
			Annotations: cr.Spec.Ingress.Annotations,
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, ing, func() error {

		// We don't accept any update for service
		if !ing.ObjectMeta.CreationTimestamp.IsZero() {
			return nil
		}

		if err := controllerutil.SetControllerReference(cr, ing, r.Scheme); err != nil {
			return err
		}

		// we only create ingress rule with backend service created by this operator.
		ingressRule := networkingv1beta1.IngressRuleValue{
			HTTP: &networkingv1beta1.HTTPIngressRuleValue{
				Paths: []networkingv1beta1.HTTPIngressPath{
					{
						Backend: networkingv1beta1.IngressBackend{
							ServiceName: cr.GetName(),
							ServicePort: intstr.FromInt(blogPort),
						},
					},
				},
			},
		}

		hosts := cr.Spec.Ingress.Hosts
		var rules []networkingv1beta1.IngressRule
		// if no one hosts defined, create rule without any spesific host.
		if len(hosts) == 0 {
			rules = append(rules, networkingv1beta1.IngressRule{
				IngressRuleValue: ingressRule,
			})
		} else {
			for _, host := range hosts {
				rules = append(rules, networkingv1beta1.IngressRule{
					Host:             host,
					IngressRuleValue: ingressRule,
				})
			}
		}

		ing.Spec.Rules = rules
		// if tls is enabled add `IngressTLS` with defined hosts.
		// NOTE: hosts is required when tls is enabled.
		if cr.Spec.Ingress.TLS.Enabled {
			ing.Spec.TLS = append(ing.Spec.TLS, networkingv1beta1.IngressTLS{
				Hosts:      hosts,
				SecretName: cr.Spec.Ingress.TLS.SecretName,
			})
		}

		return nil
	})

	r.Logger.Info("Reconciling Blog Ingress: ", ing.Name)
	return err
}

func (r *BlogReconciler) CreateOrUpdateService(cr *ghostv1alpha1.Blog) error {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    commonLabelFromCR(cr),
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
			Selector: commonLabelFromCR(cr),
			Type:     cr.Spec.ServiceType,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   "TCP",
					Port:       int32(blogPort),
					TargetPort: intstr.FromInt(blogPort),
				},
			},
		}
		return nil
	})

	r.Logger.Info("Reconciling Service: ", svc.Name)
	return err
}

func (r *BlogReconciler) CreateOrUpdatePersistentVolumeClaim(cr *ghostv1alpha1.Blog) error {
	requestStorage := make(corev1.ResourceList)
	requestStorage[corev1.ResourceStorage] = cr.Spec.Persistent.Size
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      persistentVolumeClaimNameFromCR(cr),
			Namespace: cr.GetNamespace(),
			Labels:    commonLabelFromCR(cr),
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

	r.Logger.Info("Reconciling PersistentVolumeClaim: ", pvc.Name)
	return err
}

func (r *BlogReconciler) CreateOrUpdateConfigMap(cr *ghostv1alpha1.Blog) error {
	configdata := make(map[string]string)

	cr.Spec.Config.Server.Host = "0.0.0.0"
	cr.Spec.Config.Server.Port = intstr.FromInt(blogPort)
	cr.Spec.Config.Database.Connection.Host = cr.Name + mysqlSuffix

	config, _ := json.MarshalIndent(cr.Spec.Config, "", "  ")
	configdata["config.json"] = string(config)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapNameFromCR(cr),
			Namespace: cr.GetNamespace(),
			Labels:    commonLabelFromCR(cr),
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, cm, func() error {
		if err := controllerutil.SetControllerReference(cr, cm, r.Scheme); err != nil {
			return err
		}

		cm.Data = configdata
		return nil
	})

	r.Logger.Info("Reconciling  ConfigMap: ", cm.Name)
	return err
}

func (r *BlogReconciler) CreateOrUpdateDeployment(cr *ghostv1alpha1.Blog) error {
	defaultTerminationGracePeriodSeconds := int64(30)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    commonLabelFromCR(cr),
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.Client, dep, func() error {
		if dep.ObjectMeta.CreationTimestamp.IsZero() {
			// Set label selector only when deployment has never been created
			dep.Spec.Selector = commonLabelSelectorFromCR(cr)
		}

		if err := controllerutil.SetControllerReference(cr, dep, r.Scheme); err != nil {
			return err
		}

		dep.Spec.Replicas = cr.Spec.Replicas
		dep.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: commonLabelFromCR(cr),
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "ghost",
						Image:           cr.Spec.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{
							{
								Name:          "http",
								ContainerPort: int32(blogPort),
								Protocol:      corev1.ProtocolTCP,
							},
						},
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"/bin/sh", "-c", "ln -sf /etc/ghost/config/config.json /var/lib/ghost/config.production.json"},
								},
							},
						},
						TerminationMessagePath:   "/dev/termination-log",
						TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						VolumeMounts:             r.newVolumeMountForCR(),
					},
				},
				RestartPolicy:                 corev1.RestartPolicyAlways,
				TerminationGracePeriodSeconds: &defaultTerminationGracePeriodSeconds,
				DNSPolicy:                     corev1.DNSClusterFirst,
				SecurityContext:               &corev1.PodSecurityContext{},
				SchedulerName:                 corev1.DefaultSchedulerName,
				Volumes:                       r.newVolumeForCR(cr),
			},
		}
		return nil
	})

	r.Logger.Info("Reconciling Deployment: ", dep.Name)
	return err
}

func (r *BlogReconciler) newVolumeMountForCR() []corev1.VolumeMount {
	var volumeMount []corev1.VolumeMount

	volumeMount = append(volumeMount, corev1.VolumeMount{
		Name:      "ghost-config",
		ReadOnly:  true,
		MountPath: "/etc/ghost/config",
	})
	volumeMount = append(volumeMount, corev1.VolumeMount{
		Name:      "ghost-content",
		ReadOnly:  false,
		MountPath: "/var/lib/ghost/content",
	})

	return volumeMount

}
