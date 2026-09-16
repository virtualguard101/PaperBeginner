"""Smoke test PaperBeginner FastAPI demo endpoints."""
from pathlib import Path

import httpx

base = "http://127.0.0.1:8000"
c = httpx.Client(base_url=base, timeout=120.0)

r = c.get("/health")
assert r.status_code == 200, r.text
print("health", r.json())

email = "demo@example.com"
password = "password123"
r = c.post("/api/v1/auth/register", json={"email": email, "password": password, "name": "Demo User"})
print("register", r.status_code, r.json().get("success") or r.json().get("error"))

r = c.post("/api/v1/auth/login", json={"email": email, "password": password})
assert r.status_code == 200 and r.json()["success"], r.text
token = r.json()["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}
print("login ok")

r = c.get("/api/v1/users/me", headers=headers)
assert r.json()["success"], r.text
print("me", r.json()["data"]["email"])

r = c.get("/api/v1/ccf/categories", headers=headers)
assert r.json()["success"] and len(r.json()["data"]) == 10
print("ccf", len(r.json()["data"]))

r = c.get("/api/v1/trending", headers=headers, params={"limit": 20})
assert r.json()["success"], r.text
print("trending", len(r.json()["data"]))

r = c.get("/api/v1/trending/reports", headers=headers)
assert r.json()["success"], r.text
print("reports", len(r.json()["data"]))

pdf = Path("frontend/public/sample-attention.pdf")
files = {"file": ("sample-attention.pdf", pdf.read_bytes(), "application/pdf")}
data = {"title": "Attention Is All You Need Demo"}
r = c.post("/api/v1/papers/upload", headers=headers, files=files, data=data)
assert r.status_code == 201 and r.json()["success"], r.text
paper_id = r.json()["data"]["id"]
print("upload", paper_id)

r = c.post(
    f"/api/v1/papers/{paper_id}/analyze",
    headers=headers,
    json={"types": ["summary"]},
)
assert r.json()["success"], r.text
print("analyze", r.json()["data"]["status"], "analyses", len(r.json()["data"].get("analyses") or []))

r = c.post(
    "/api/v1/learning/generate",
    headers=headers,
    json={"category_id": 8, "difficulty": "beginner"},
)
assert r.status_code == 201 and r.json()["success"], r.text
print("learning", r.json()["data"]["title"])

r = c.post(
    "/api/v1/reviews/generate",
    headers=headers,
    json={"title": "Transformer Demo Review", "paper_ids": [paper_id], "style": "academic"},
)
assert r.status_code == 201 and r.json()["success"], r.text
review_id = r.json()["data"]["id"]
print("review", review_id)

r = c.post(f"/api/v1/reviews/{review_id}/score", headers=headers, json={})
assert r.json()["success"] and r.json()["data"]["status"] == "scored", r.text
print("score", r.json()["data"]["score"]["overall_score"])

print("ALL SMOKE TESTS PASSED")
