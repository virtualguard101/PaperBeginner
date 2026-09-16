.PHONY: help demo demo-backend demo-frontend install docker-up docker-down smoke

help:
	@echo "PaperBeginner demo targets"
	@echo "  make install         Install frontend + backend deps"
	@echo "  make demo-backend    Run FastAPI on :8000"
	@echo "  make demo-frontend   Run Vite on :5173"
	@echo "  make smoke           API smoke test (server must be running)"
	@echo "  make docker-up       docker compose up --build"
	@echo "  make docker-down     docker compose down"

install:
	cd backend && python -m pip install -r requirements.txt
	cd frontend && npm install

demo-backend:
	cd backend && python -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

demo-frontend:
	cd frontend && npm run dev

smoke:
	python scripts/smoke_test.py

docker-up:
	docker compose up --build

docker-down:
	docker compose down
