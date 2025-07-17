package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

var (
	create   bool
	update   bool
	delete   bool
	replicas int
	image    string
)

var k8sdeploy = &cobra.Command{
	Use:   "k8sdeploy",
	Short: "Manage Kubernetes resources for Book Project",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Kubernetes resource management")

		// Validate flags
		if (create && (update || delete)) || (update && delete) || (!create && !update && !delete) {
			log.Fatalf("Please specify exactly one of --create, --update, or --delete")
		}

		// Load kubeconfig
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Failed to build kubeconfig: %v", err)
		}

		// Create clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Failed to create clientset: %v", err)
		}

		// Execute operation
		if create {
			log.Println("Creating Kubernetes resources")
			createResources(clientset)
		} else if update {
			log.Println("Updating Deployment")
			updateDeployment(clientset)
		} else if delete {
			log.Println("Deleting Kubernetes resources")
			deleteResources(clientset)
		}
	},
}

func init() {
	rootCmd.AddCommand(k8sdeploy)
	k8sdeploy.PersistentFlags().BoolVarP(&create, "create", "c", false, "Create the book-project ConfigMap, Deployment, and Service")
	k8sdeploy.PersistentFlags().BoolVarP(&update, "update", "u", false, "Update the book-project Deployment (replicas or image)")
	k8sdeploy.PersistentFlags().BoolVarP(&delete, "delete", "d", false, "Delete the book-project ConfigMap, Deployment, and Service")
	k8sdeploy.PersistentFlags().IntVarP(&replicas, "replicas", "r", 1, "Number of replicas for the Deployment (used with --update)")
	k8sdeploy.PersistentFlags().StringVarP(&image, "image", "i", "yourdockerhubusername/book-project:latest", "Docker image for the Deployment (used with --update)")
}

func createResources(clientset *kubernetes.Clientset) {
	// Create ConfigMap
	log.Println("Creating ConfigMap 'book-project-config'")
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "book-project-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"jwt-secret": "bolaJabeNah",
		},
	}
	_, err := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Fatalf("Failed to create ConfigMap: %v", err)
	}
	log.Println("Created ConfigMap 'book-project-config'")

	// Create Deployment
	log.Println("Creating Deployment 'book-project-deployment'")
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "book-project-deployment",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "book-project",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "book-project",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "book-project",
							Image: "yourdockerhubusername/book-project:latest",
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080},
							},
							Env: []corev1.EnvVar{
								{
									Name: "JWT_SECRET",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{Name: "book-project-config"},
											Key:                  "jwt-secret",
										},
									},
								},
							},
							Command: []string{"./main", "startProject", "--auth=true"},
						},
					},
				},
			},
		},
	}
	_, err = clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Fatalf("Failed to create Deployment: %v", err)
	}
	log.Println("Created Deployment 'book-project-deployment'")

	// Create Service
	log.Println("Creating Service 'book-project-service'")
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "book-project-service",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "book-project",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(8080),
					NodePort:   30081,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	_, err = clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Fatalf("Failed to create Service: %v", err)
	}
	log.Println("Created Service 'book-project-service'")
}

func updateDeployment(clientset *kubernetes.Clientset) {
	log.Println("Updating Deployment 'book-project-deployment'")
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deploymentsClient := clientset.AppsV1().Deployments("default")
		deployment, err := deploymentsClient.Get(context.TODO(), "book-project-deployment", metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Failed to get Deployment: %v", err)
		}
		deployment.Spec.Replicas = int32Ptr(int32(replicas))
		deployment.Spec.Template.Spec.Containers[0].Image = image
		_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
		return err
	})
	if err != nil {
		log.Fatalf("Failed to update Deployment: %v", err)
	}
	log.Println("Updated Deployment 'book-project-deployment'")
}

func deleteResources(clientset *kubernetes.Clientset) {
	log.Println("Deleting Service 'book-project-service'")
	err := clientset.CoreV1().Services("default").Delete(context.TODO(), "book-project-service", metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Fatalf("Failed to delete Service: %v", err)
	}
	log.Println("Deleted Service 'book-project-service'")

	log.Println("Deleting Deployment 'book-project-deployment'")
	err = clientset.AppsV1().Deployments("default").Delete(context.TODO(), "book-project-deployment", metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Fatalf("Failed to delete Deployment: %v", err)
	}
	log.Println("Deleted Deployment 'book-project-deployment'")

	log.Println("Deleting ConfigMap 'book-project-config'")
	err = clientset.CoreV1().ConfigMaps("default").Delete(context.TODO(), "book-project-config", metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Fatalf("Failed to delete ConfigMap: %v", err)
	}
	log.Println("Deleted ConfigMap 'book-project-config'")
}

func int32Ptr(i int32) *int32 { return &i }