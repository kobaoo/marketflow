dc_up:
	docker-compose up --build
dc_down:
	docker-compose down
exchange:
	docker load -i exchange_tars/exchange1_amd64.tar
	docker load -i exchange_tars/exchange2_amd64.tar
	docker load -i exchange_tars/exchange3_amd64.tar
run_exchanges:
	docker run -p 40101:40101 --name exchange1 -d exchange1:latest
	docker run -p 40102:40102 --name exchange2 -d exchange2:latest
	docker run -p 40103:40103 --name exchange3 -d exchange3:latest
rerun_exchanges:
	docker run exchange1
	docker run exchange2
	docker run exchange3
stop_exchanges:
	docker stop exchange1
	docker stop exchange2
	docker stop exchange3
rm_exchanges:
	docker rm exchange1
	docker rm exchange2
	docker rm exchange3