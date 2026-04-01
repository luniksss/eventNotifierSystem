## участники
1) api-gateway
2) task-processor
3) notifier

## пипелине
1) api-gateway process requests with /task
2) api-gateway creates task, store it to database
3) api-gateway creates event and publish it to first channel
4) task-processor(first channel's consumer) gets event and validates it
5) if everything is ok, task-processor creates event and publish it to another channel
6) notifier(second channel's consumer) gets event and send data to mock users

## minikube команды потому что я глупий-глупий с плохой память
1) minikube start --cpus=4 --memory=8192 --driver=docker
2) kubectl get nodes
3) eval $(minikube docker-env)
4) cd api-gateway
   docker build -t eventnotifiersystem-api-gateway:latest .
   cd ../task-processor
   docker build -t eventnotifiersystem-task-processor:latest .
   cd ../notifier
   docker build -t eventnotifiersystem-notifier:latest .
   cd ..
5) kubectl create namespace event-notification
6) kubectl create configmap mysql-init-sql --from-file=scripts/init.sql -n event-notification
7) kubectl create secret generic app-secrets -n event-notification \
   --from-literal=db-dsn='test:1234@tcp(notifier-mysql:3306)/tasksdb?parseTime=true' \
   --from-literal=rabbitmq-url='amqp://guest:guest@notifier-rabbbitmq:5672/'
8) kubectl create configmap app-config -n event-notification \
   --from-literal=tasks-exchange='tasks' \
   --from-literal=tasks-routing-key='task.created' \
   --from-literal=notifications-exchange='notifications'
9) kubectl apply -f k8s/mysql/
10) kubectl apply -f k8s/rabbitmq/
11) kubectl apply -f k8s/api-gateway/
12) kubectl apply -f k8s/task-processor/
13) kubectl apply -f k8s/notifier/
14) kubectl get pods -n event-notification -w
15) kubectl describe pod api-gateway-... -n event-notification
16) kubectl port-forward -n event-notification deployment/api-gateway 8080:8080
17) kubectl logs -n event-notification deployment/api-gateway

## примерчик запросика
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{"title":"Test","description":"Desc","email":"test@test.com","phone":"88009993434"}'