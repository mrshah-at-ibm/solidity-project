package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

func SetupClient() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfigMap(namespace, name string) (map[string]string, error) {
	var err error
	if clientset == nil {
		err = SetupClient()
		if err != nil {
			return nil, err
		}
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm.Data, nil
}

func WriteConfigMap(namespace, name string, data map[string]string) error {
	var err error
	if clientset == nil {
		err = SetupClient()
		if err != nil {
			return err
		}
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	cm.Data = data
	_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ReadSecret(namespace, name string) (map[string]string, error) {
	var err error
	if clientset == nil {
		err = SetupClient()
		if err != nil {
			return nil, err
		}
	}

	s, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return s.StringData, nil

}

func WriteSecret(namespace, name string, data map[string]string) error {
	var err error
	if clientset == nil {
		err = SetupClient()
		if err != nil {
			return err
		}
	}

	s, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			s1 := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				StringData: data,
			}
			_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), s1, metav1.CreateOptions{})
			if err != nil {
				return err
			}

			return nil

		}
		return err
	}

	s.StringData = data
	_, err = clientset.CoreV1().Secrets(namespace).Update(context.TODO(), s, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
