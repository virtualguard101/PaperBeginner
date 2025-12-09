-- Drop triggers
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TRIGGER IF EXISTS update_learning_paths_updated_at ON learning_paths;
DROP TRIGGER IF EXISTS update_trending_items_updated_at ON trending_items;
DROP TRIGGER IF EXISTS update_papers_updated_at ON papers;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS review_papers;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS learning_paths;
DROP TABLE IF EXISTS trend_reports;
DROP TABLE IF EXISTS trending_items;
DROP TABLE IF EXISTS paper_analyses;
DROP TABLE IF EXISTS papers;
DROP TABLE IF EXISTS ccf_categories;
DROP TABLE IF EXISTS users;

-- Note: We don't drop the extensions as they might be used by other databases

