run_order_service:
	cd ./services/order_service && docker-compose up -d

start: run_order_service
	ls
