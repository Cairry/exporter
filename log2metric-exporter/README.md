# log2metric

日志转metric

log格式：
```
{
    "@timestamp": "2023-11-27T07:17:32.009Z",
    "k8s_container_name": "buxybox",
    "k8s_pod_namespace": "js-design-prod",
    "message": "{\"timestamp\":\"2023-11-23T06:34:56.980867Z\",\"level\":\"INFO\"}"
}
```
metric格式：
```
l2m_level_info{level="ERROR",namespace="ztest",service="wood"} 1
l2m_level_info{level="INFO",namespace="ztest",service="wood"} 1
l2m_level_info{level="DEBUG",namespace="ztest",service="wood"} 1
l2m_level_info{level="Unknown",namespace="ztest",service="wood"} 1
```
