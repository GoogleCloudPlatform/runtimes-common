#!/usr/bin/python

GKE_TEMPLATE = '''apiVersion: v1
kind: Service
metadata:
  name: {service_name}
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: {service_name}-app
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {service_name}
spec:
  replicas: 1
  template:
    metadata:
      name: {service_name}
      labels:
        app: {service_name}-app
    spec:
      containers:
      - name: {service_name}
        image: {test_image}
        resources:
          requests:
            cpu: 100m
            memory: "512Mi"
          limits:
            cpu: 200m
            memory: "1024Mi"
        env:
        - name: HEAP_SIZE_RATIO
          value: "50"
        ports:
        - containerPort: 8080

'''
