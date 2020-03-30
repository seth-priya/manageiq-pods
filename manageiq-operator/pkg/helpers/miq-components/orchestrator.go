package miqtools

import (
	miqv1alpha1 "github.com/manageiq-operator/pkg/apis/manageiq/v1alpha1"
	randstring "github.com/manageiq-operator/pkg/helpers/randstring"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewOrchestratorDeployment(cr *miqv1alpha1.Manageiq) *appsv1.Deployment {
	DeploymentLabels := map[string]string{
		"app": cr.Spec.AppName,
	}

	PodLabels := map[string]string{
		"name": "orchestrator",
		"app":  cr.Spec.AppName,
	}

	var RepNum int32 = 1
	var termSecs int64 = 90
	memLimit, _ := resource.ParseQuantity(cr.Spec.OrchestratorMemLimit)
	memReq, _ := resource.ParseQuantity(cr.Spec.OrchestratorMemReq)
	cpuReq, _ := resource.ParseQuantity(cr.Spec.OrchestratorCPUReq)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orchestrator",
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    DeploymentLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
			Replicas: &RepNum,
			Selector: &metav1.LabelSelector{
				MatchLabels: PodLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "orchestrator",
					Labels: PodLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "orchestrator",
							Image: cr.Spec.OrchestratorImageNamespace + "/" + cr.Spec.OrchestratorImageName + ":" + cr.Spec.OrchestratorImageTag,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"pidof", "MIQ Server"},
									},
								},
								InitialDelaySeconds: 480,
								TimeoutSeconds:      3,
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ALLOW_INSECURE_SESSION",
									Value: "true",
								},
								corev1.EnvVar{
									Name:  "APP_NAME",
									Value: cr.Spec.AppName,
								},
								corev1.EnvVar{
									Name: "APPLICATION_ADMIN_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "app-secrets"},
											Key:                  "admin-password",
										},
									},
								},
								corev1.EnvVar{
									Name:  "GUID",
									Value: randstring.GenerateGUID(),
								},

								corev1.EnvVar{
									Name:  "DATABASE_REGION",
									Value: "0",
								},
								corev1.EnvVar{
									Name: "DATABASE_HOSTNAME",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "postgresql-secrets"},
											Key:                  "hostname",
										},
									},
								},
								corev1.EnvVar{
									Name: "DATABASE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "postgresql-secrets"},
											Key:                  "dbname",
										},
									},
								},
								corev1.EnvVar{
									Name: "DATABASE_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "postgresql-secrets"},
											Key:                  "password",
										},
									},
								},
								corev1.EnvVar{
									Name:  "DATABASE_PORT",
									Value: cr.Spec.DatabasePort,
								},
								corev1.EnvVar{
									Name: "DATABASE_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "postgresql-secrets"},
											Key:                  "username",
										},
									},
								},
								corev1.EnvVar{
									Name: "ENCRYPTION_KEY",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "app-secrets"},
											Key:                  "encryption-key",
										},
									},
								},
								corev1.EnvVar{
									Name:  "CONTAINER_IMAGE_NAMESPACE",
									Value: cr.Spec.OrchestratorImageNamespace,
								},
								corev1.EnvVar{
									Name:  "IMAGE_PULL_SECRET",
									Value: "",
								},
							},

							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"memory": memLimit,
								},
								Requests: corev1.ResourceList{
									"memory": memReq,
									"cpu":    cpuReq,
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						corev1.LocalObjectReference{Name: ""},
					},
					TerminationGracePeriodSeconds: &termSecs,

					ServiceAccountName: cr.Spec.AppName + "-orchestrator",
				},
			},
		},
	}
}
