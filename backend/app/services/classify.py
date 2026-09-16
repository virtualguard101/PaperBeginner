KEYWORDS = [
    (8, ["ai", "ml", "llm", "deep-learning", "transformer", "nlp", "cv", "pytorch", "tensorflow", "agent"]),
    (3, ["security", "crypto", "cve", "auth", "privacy", "malware"]),
    (2, ["network", "tcp", "http", "quic", "p2p", "cdn"]),
    (4, ["compiler", "os", "kernel", "rust", "programming-language", "devops"]),
    (5, ["database", "sql", "search", "retrieval", "vector"]),
    (1, ["gpu", "cuda", "distributed", "storage", "kubernetes"]),
    (7, ["graphics", "render", "vision", "image", "video", "diffusion"]),
    (9, ["hci", "ux", "ar", "vr", "ubiquitous"]),
    (6, ["algorithm", "complexity", "theory", "proof"]),
]


def classify_text(*parts: str) -> int:
    blob = " ".join(p or "" for p in parts).lower()
    best_id = 10
    best_hits = 0
    for cat_id, words in KEYWORDS:
        hits = sum(1 for w in words if w in blob)
        if hits > best_hits:
            best_hits = hits
            best_id = cat_id
    return best_id if best_hits else 10
