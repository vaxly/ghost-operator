package controllers

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import ghostv1alpha1 "github.com/vaxly/ghost-operator/api/v1alpha1"

const (
	blogPort    = 2368
	mysqlPort   = 3306
	mysqlSuffix = "-mysql"
)

func commonLabelFromCR(cr *ghostv1alpha1.Blog) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":     cr.Name,
		"app.kubernetes.io/instance": cr.Name,
	}
}
func commonLabelSelectorFromCR(cr *ghostv1alpha1.Blog) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: commonLabelFromCR(cr),
	}
}

func mysqlLabelFromCR(cr *ghostv1alpha1.Blog) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":     cr.Name + "-mysql",
		"app.kubernetes.io/instance": cr.Name + "-mysql",
	}
}
func mysqlLabelSelectorFromCR(cr *ghostv1alpha1.Blog) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: mysqlLabelFromCR(cr),
	}
}

func configMapNameFromCR(cr *ghostv1alpha1.Blog) string { return cr.GetName() + "-ghost-config" }

func persistentVolumeClaimNameFromCR(cr *ghostv1alpha1.Blog) string {
	return cr.GetName() + "-ghost-content-pvc"
}
func mysqlPersistentVolumeClaimNameFromCR(cr *ghostv1alpha1.Blog) string {
	return cr.GetName() + "-ghost-mysql-pvc"
}
