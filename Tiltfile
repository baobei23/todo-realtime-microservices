# Load the restart_process extension

load('ext://restart_process', 'docker_build_with_restart')
### K8s Config ###

# Uncomment to use secrets
k8s_yaml('./infra/development/k8s/secrets.yaml')
k8s_yaml('./infra/development/k8s/app-config.yaml')

### End of K8s Config ###

### RabbitMQ ###
k8s_yaml('./infra/development/k8s/rabbitmq-deployment.yaml')
k8s_resource('rabbitmq', port_forwards=['5672', '15672'], labels='tooling')
### End RabbitMQ ###

### PostgreSQL ###
k8s_yaml('./infra/development/k8s/postgres-deployment.yaml')
k8s_resource('postgres', port_forwards=['5432'], labels='tooling')
### End PostgreSQL ###

### API Gateway ###

gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")


docker_build_with_restart(
  'todo-realtime-microservices/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile', 'rabbitmq'], labels="services")
### End of API Gateway ###

### Todo Service ###

todo_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/todo-service ./services/todo-service/cmd/main.go'
if os.name == 'nt':
 todo_compile_cmd = './infra/development/docker/todo-build.bat'

local_resource(
  'todo-service-compile',
  todo_compile_cmd,
  deps=['./services/todo-service', './shared'], labels="compiles")

docker_build_with_restart(
  'todo-realtime-microservices/todo-service',
  '.',
  entrypoint=['/app/build/todo-service'],
  dockerfile='./infra/development/docker/todo-service.Dockerfile',
  only=[
    './build/todo-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/todo-service-deployment.yaml')
k8s_resource('todo-service', resource_deps=['todo-service-compile', 'rabbitmq'], labels="services")

### End of Todo Service ###

### Realtime Service ###

realtime_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/realtime-service ./services/realtime-service/cmd/main.go'
if os.name == 'nt':
 realtime_compile_cmd = './infra/development/docker/realtime-build.bat'

local_resource(
  'realtime-service-compile',
  realtime_compile_cmd,
  deps=['./services/realtime-service', './shared'], labels="compiles")

docker_build_with_restart(
  'todo-realtime-microservices/realtime-service',
  '.',
  entrypoint=['/app/build/realtime-service'],
  dockerfile='./infra/development/docker/realtime-service.Dockerfile',
  only=[
    './build/realtime-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/realtime-service-deployment.yaml')
k8s_resource('realtime-service', resource_deps=['realtime-service-compile', 'rabbitmq'], labels="services")

### End of Realtime Service ###