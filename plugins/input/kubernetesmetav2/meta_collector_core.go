package kubernetesmetav2

import (
	"encoding/json"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/alibaba/ilogtail/pkg/helper/k8smeta"
	"github.com/alibaba/ilogtail/pkg/models"
)

func (m *metaCollector) processPodEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	result := []models.PipelineEvent{}
	if obj, ok := data.Raw.(*v1.Pod); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, obj.Namespace, obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("api_version", obj.APIVersion)
		log.Contents.Add("namespace", obj.Namespace)
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		log.Contents.Add("annotations", m.processEntityJSONObject(obj.Annotations))
		log.Contents.Add("status", string(obj.Status.Phase))
		log.Contents.Add("instance_ip", obj.Status.PodIP)
		containerInfos := []map[string]string{}
		for _, container := range obj.Spec.Containers {
			containerInfo := map[string]string{
				"name":  container.Name,
				"image": container.Image,
			}
			containerInfos = append(containerInfos, containerInfo)
		}
		log.Contents.Add("containers", m.processEntityJSONArray(containerInfos))
		result = append(result, log)

		// container
		if m.serviceK8sMeta.Container {
			for _, container := range obj.Spec.Containers {
				containerLog := &models.Log{}
				containerLog.Contents = models.NewLogContents()
				containerLog.Timestamp = log.Timestamp

				containerLog.Contents.Add(entityDomainFieldName, m.serviceK8sMeta.domain)
				containerLog.Contents.Add(entityTypeFieldName, m.genEntityTypeKey(containerTypeName))
				containerLog.Contents.Add(entityIDFieldName, m.genKey(containerTypeName, obj.Namespace, obj.Name+container.Name))
				containerLog.Contents.Add(entityMethodFieldName, method)

				containerLog.Contents.Add(entityFirstObservedTimeFieldName, strconv.FormatInt(data.FirstObservedTime, 10))
				containerLog.Contents.Add(entityLastObservedTimeFieldName, strconv.FormatInt(data.LastObservedTime, 10))
				containerLog.Contents.Add(entityKeepAliveSecondsFieldName, strconv.FormatInt(int64(m.serviceK8sMeta.Interval*2), 10))
				containerLog.Contents.Add(entityCategoryFieldName, defaultEntityCategory)

				// common custom fields
				containerLog.Contents.Add(entityClusterIDFieldName, m.serviceK8sMeta.clusterID)
				containerLog.Contents.Add(entityNameFieldName, container.Name)

				// custom fields
				containerLog.Contents.Add("pod_name", obj.Name)
				containerLog.Contents.Add("pod_namespace", obj.Namespace)
				containerLog.Contents.Add("image", container.Image)
				if container.Resources.Requests != nil && container.Resources.Requests.Cpu() != nil {
					containerLog.Contents.Add("cpu_request", container.Resources.Requests.Cpu().String())
				} else {
					containerLog.Contents.Add("cpu_request", "")
				}
				if container.Resources.Limits != nil && container.Resources.Limits.Cpu() != nil {
					containerLog.Contents.Add("cpu_limit", container.Resources.Limits.Cpu().String())
				} else {
					containerLog.Contents.Add("cpu_limit", "")
				}
				if container.Resources.Requests != nil && container.Resources.Requests.Memory() != nil {
					containerLog.Contents.Add("memory_request", container.Resources.Requests.Memory().String())
				} else {
					containerLog.Contents.Add("memory_request", "")
				}
				if container.Resources.Limits != nil && container.Resources.Limits.Memory() != nil {
					containerLog.Contents.Add("memory_limit", container.Resources.Limits.Memory().String())
				} else {
					containerLog.Contents.Add("memory_limit", "")
				}
				ports := make([]int32, 0)
				if container.Ports != nil {
					for _, port := range container.Ports {
						ports = append(ports, port.ContainerPort)
					}
				}
				portsStr, _ := json.Marshal(ports)
				if len(ports) == 0 {
					portsStr = []byte("[]")
				}
				containerLog.Contents.Add("container_ports", string(portsStr))
				volumes := make([]map[string]string, 0)
				if container.VolumeMounts != nil {
					for _, volume := range container.VolumeMounts {
						volumeInfo := map[string]string{
							"volumeMountName": volume.Name,
							"volumeMountPath": volume.MountPath,
						}
						volumes = append(volumes, volumeInfo)
					}
				}
				containerLog.Contents.Add("volumes", m.processEntityJSONArray(volumes))
				result = append(result, containerLog)
			}
		}
		return result
	}
	return nil
}

func (m *metaCollector) processNodeEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*v1.Node); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, "", obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		log.Contents.Add("annotations", m.processEntityJSONObject(obj.Annotations))
		status := []map[string]string{}
		for _, condition := range obj.Status.Conditions {
			conditionInfo := map[string]string{
				"type":   string(condition.Type),
				"status": string(condition.Status),
			}
			status = append(status, conditionInfo)
		}
		log.Contents.Add("status", m.processEntityJSONArray(status))
		for _, addr := range obj.Status.Addresses {
			if addr.Type == v1.NodeInternalIP {
				log.Contents.Add("internal_ip", addr.Address)
			} else if addr.Type == v1.NodeHostName {
				log.Contents.Add("host_name", addr.Address)
			}
		}
		capacityStr, _ := json.Marshal(obj.Status.Capacity)
		log.Contents.Add("capacity", string(capacityStr))
		allocatableStr, _ := json.Marshal(obj.Status.Allocatable)
		log.Contents.Add("allocatable", string(allocatableStr))
		addressStr, _ := json.Marshal(obj.Status.Addresses)
		log.Contents.Add("addresses", string(addressStr))
		log.Contents.Add("provider_id", obj.Spec.ProviderID)
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processServiceEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*v1.Service); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, obj.Namespace, obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("api_version", obj.APIVersion)
		log.Contents.Add("namespace", obj.Namespace)
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		log.Contents.Add("annotations", m.processEntityJSONObject(obj.Annotations))
		log.Contents.Add("selector", m.processEntityJSONObject(obj.Spec.Selector))
		log.Contents.Add("type", string(obj.Spec.Type))
		log.Contents.Add("cluster_ip", obj.Spec.ClusterIP)
		ports := make([]map[string]string, 0)
		for _, port := range obj.Spec.Ports {
			targetPort := ""
			if port.TargetPort.Type == 0 { // IntOrString type, 0 means int
				targetPort = strconv.FormatInt(int64(port.TargetPort.IntVal), 10)
			} else {
				targetPort = port.TargetPort.StrVal
			}
			portInfo := map[string]string{
				"port":       strconv.FormatInt(int64(port.Port), 10),
				"targetPort": targetPort,
				"protocol":   string(port.Protocol),
			}
			ports = append(ports, portInfo)
		}
		log.Contents.Add("ports", m.processEntityJSONArray(ports))
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processConfigMapEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*v1.ConfigMap); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, obj.Namespace, obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("api_version", obj.APIVersion)
		log.Contents.Add("namespace", obj.Namespace)
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		log.Contents.Add("annotations", m.processEntityJSONObject(obj.Annotations))
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processNamespaceEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*v1.Namespace); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, "", obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("api_version", obj.APIVersion)
		log.Contents.Add("kind", obj.Kind)
		log.Contents.Add("name", obj.Name)
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPersistentVolumeEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*v1.PersistentVolume); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, obj.Namespace, obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("api_version", obj.APIVersion)
		log.Contents.Add("namespace", obj.Namespace)
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		log.Contents.Add("annotations", m.processEntityJSONObject(obj.Annotations))
		log.Contents.Add("status", string(obj.Status.Phase))
		log.Contents.Add("storage_class_name", obj.Spec.StorageClassName)
		log.Contents.Add("persistent_volume_reclaim_policy", string(obj.Spec.PersistentVolumeReclaimPolicy))
		if obj.Spec.VolumeMode != nil {
			log.Contents.Add("volume_mode", string(*obj.Spec.VolumeMode))
		} else {
			log.Contents.Add("volume_mode", "")
		}
		if obj.Spec.Capacity != nil && obj.Spec.Capacity.Storage() != nil {
			log.Contents.Add("capacity", obj.Spec.Capacity.Storage().String())
		} else {
			log.Contents.Add("capacity", "")
		}
		if obj.Spec.CSI != nil {
			log.Contents.Add("fsType", obj.Spec.CSI.FSType)
		} else {
			log.Contents.Add("fsType", "")
		}
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPersistentVolumeClaimEntity(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*v1.PersistentVolumeClaim); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		log.Timestamp = uint64(time.Now().Unix())
		m.processEntityCommonPart(log.Contents, obj.Kind, obj.Namespace, obj.Name, method, data.FirstObservedTime, data.LastObservedTime, obj.CreationTimestamp)

		// custom fields
		log.Contents.Add("api_version", obj.APIVersion)
		log.Contents.Add("namespace", obj.Namespace)
		log.Contents.Add("labels", m.processEntityJSONObject(obj.Labels))
		log.Contents.Add("annotations", m.processEntityJSONObject(obj.Annotations))
		log.Contents.Add("status", string(obj.Status.Phase))
		log.Contents.Add("storeage_requests", obj.Spec.Resources.Requests.Storage().String())
		if obj.Spec.StorageClassName != nil {
			log.Contents.Add("storage_class_name", *obj.Spec.StorageClassName)
		} else {
			log.Contents.Add("storage_class_name", "")
		}
		log.Contents.Add("volume_name", obj.Spec.VolumeName)
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPodNodeLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PodNode); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Node.Kind, "", obj.Node.Name, obj.Pod.Kind, obj.Pod.Namespace, obj.Pod.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Node2Pod)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPodPVCLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PodPersistentVolumeClaim); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Pod.Kind, obj.Pod.Namespace, obj.Pod.Name, obj.PersistentVolumeClaim.Kind, obj.PersistentVolumeClaim.Namespace, obj.PersistentVolumeClaim.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Pod2PersistentVolumeClaim)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPodConfigMapLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PodConfigMap); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Pod.Kind, obj.Pod.Namespace, obj.Pod.Name, obj.ConfigMap.Kind, obj.ConfigMap.Namespace, obj.ConfigMap.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Pod2ConfigMap)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPodServiceLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PodService); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Service.Kind, obj.Service.Namespace, obj.Service.Name, obj.Pod.Kind, obj.Pod.Namespace, obj.Pod.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Service2Pod)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPodContainerLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PodContainer); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Pod.Kind, obj.Pod.Namespace, obj.Pod.Name, "container", obj.Pod.Namespace, obj.Pod.Name+obj.Container.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Pod2Container)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPodNamespaceLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PodNamespace); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Namespace.Kind, obj.Namespace.Namespace, obj.Namespace.Name, obj.Pod.Kind, obj.Pod.Namespace, obj.Pod.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Namespace2Pod)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processServiceNamespaceLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.ServiceNamespace); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Namespace.Kind, obj.Namespace.Namespace, obj.Namespace.Name, obj.Service.Kind, obj.Service.Namespace, obj.Service.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Namespace2Service)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processConfigMapNamespaceLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.ConfigMapNamespace); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Namespace.Kind, obj.Namespace.Namespace, obj.Namespace.Name, obj.ConfigMap.Kind, obj.ConfigMap.Namespace, obj.ConfigMap.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Namespace2Configmap)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}

func (m *metaCollector) processPVCNamespaceLink(data *k8smeta.ObjectWrapper, method string) []models.PipelineEvent {
	if obj, ok := data.Raw.(*k8smeta.PersistentVolumeClaimNamespace); ok {
		log := &models.Log{}
		log.Contents = models.NewLogContents()
		m.processEntityLinkCommonPart(log.Contents, obj.Namespace.Kind, obj.Namespace.Namespace, obj.Namespace.Name, obj.PersistentVolumeClaim.Kind, obj.PersistentVolumeClaim.Namespace, obj.PersistentVolumeClaim.Name, method, data.FirstObservedTime, data.LastObservedTime)
		log.Contents.Add(entityLinkRelationTypeFieldName, m.serviceK8sMeta.Namespace2PersistentVolumeClaim)
		log.Timestamp = uint64(time.Now().Unix())
		return []models.PipelineEvent{log}
	}
	return nil
}
