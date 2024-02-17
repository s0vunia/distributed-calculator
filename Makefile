# Makefile

# Переменная для хранения количества экземпляров службы agent
AGENTS ?=  1

# Цель для запуска всех сервисов с указанным количеством экземпляров службы agent
up:
	@docker-compose up --scale agent=$(AGENTS) -d --no-recreate

build:
	@docker-compose up --scale agent=$(AGENTS) -d --no-recreate --build
# Цель для остановки всех сервисов
down:
	@docker-compose down

restart:
	@docker-compose restart

rebuild:
	$(MAKE) down && $(MAKE) build

clean:
	@docker-compose down --rmi all --volumes

logs:
	@docker-compose logs -f

scale:
	@docker-compose up --scale agent=$(AGENTS) -d --no-recreate
