-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(500),
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    preferences JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);

-- CCF Categories table
CREATE TABLE IF NOT EXISTS ccf_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    name_en VARCHAR(200) NOT NULL
);

-- Insert CCF categories
INSERT INTO ccf_categories (id, name, name_en) VALUES
(1, '计算机体系结构/并行与分布计算/存储系统', 'Computer Architecture/Parallel and Distributed Computing/Storage Systems'),
(2, '计算机网络', 'Computer Networks'),
(3, '网络与信息安全', 'Network and Information Security'),
(4, '软件工程/系统软件/程序设计语言', 'Software Engineering/System Software/Programming Languages'),
(5, '数据库/数据挖掘/内容检索', 'Database/Data Mining/Content Retrieval'),
(6, '计算机科学理论', 'Computer Science Theory'),
(7, '计算机图形学与多媒体', 'Computer Graphics and Multimedia'),
(8, '人工智能', 'Artificial Intelligence'),
(9, '人机交互与普适计算', 'Human-Computer Interaction and Ubiquitous Computing'),
(10, '交叉/综合/新兴', 'Interdisciplinary/Comprehensive/Emerging')
ON CONFLICT (id) DO NOTHING;

-- Papers table
CREATE TABLE IF NOT EXISTS papers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    authors TEXT[],
    abstract TEXT,
    keywords TEXT[],
    published_at TIMESTAMP WITH TIME ZONE,
    venue VARCHAR(300),
    doi VARCHAR(100),
    source_url VARCHAR(500),
    file_path VARCHAR(500),
    file_size BIGINT,
    ccf_category_id INTEGER REFERENCES ccf_categories(id),
    ccf_rank VARCHAR(10),
    status VARCHAR(20) DEFAULT 'pending',
    embedding vector(1536),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_papers_user_id ON papers(user_id);
CREATE INDEX idx_papers_doi ON papers(doi);
CREATE INDEX idx_papers_ccf_category ON papers(ccf_category_id);
CREATE INDEX idx_papers_embedding ON papers USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Paper analyses table
CREATE TABLE IF NOT EXISTS paper_analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    paper_id UUID NOT NULL REFERENCES papers(id) ON DELETE CASCADE,
    analysis_type VARCHAR(50) NOT NULL,
    content TEXT,
    structured_data JSONB,
    llm_provider VARCHAR(50),
    llm_model VARCHAR(100),
    tokens_used INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_paper_analyses_paper_id ON paper_analyses(paper_id);

-- Trending items table
CREATE TABLE IF NOT EXISTS trending_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source VARCHAR(20) NOT NULL,
    source_id VARCHAR(200),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    url VARCHAR(500),
    stars INTEGER DEFAULT 0,
    forks INTEGER DEFAULT 0,
    language VARCHAR(50),
    topics TEXT[],
    ccf_category_id INTEGER REFERENCES ccf_categories(id),
    trend_score FLOAT DEFAULT 0,
    crawled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_trending_source ON trending_items(source);
CREATE INDEX idx_trending_source_id ON trending_items(source, source_id);
CREATE INDEX idx_trending_category ON trending_items(ccf_category_id);
CREATE UNIQUE INDEX idx_trending_unique ON trending_items(source, source_id);

-- Trend reports table
CREATE TABLE IF NOT EXISTS trend_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    period VARCHAR(20) NOT NULL,
    category_id INTEGER REFERENCES ccf_categories(id),
    summary TEXT,
    highlights JSONB,
    generated_by VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_trend_reports_category ON trend_reports(category_id);
CREATE INDEX idx_trend_reports_period ON trend_reports(period);

-- Learning paths table
CREATE TABLE IF NOT EXISTS learning_paths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    category_id INTEGER REFERENCES ccf_categories(id),
    title VARCHAR(300) NOT NULL,
    description TEXT,
    stages JSONB,
    prerequisites TEXT[],
    estimated_time VARCHAR(50),
    difficulty VARCHAR(20),
    generated_by VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_learning_paths_user_id ON learning_paths(user_id);
CREATE INDEX idx_learning_paths_category ON learning_paths(category_id);

-- Reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    abstract TEXT,
    content TEXT,
    paper_ids UUID[],
    category_id INTEGER REFERENCES ccf_categories(id),
    status VARCHAR(20) DEFAULT 'draft',
    score JSONB,
    generated_by VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_category ON reviews(category_id);

-- Review papers junction table
CREATE TABLE IF NOT EXISTS review_papers (
    review_id UUID REFERENCES reviews(id) ON DELETE CASCADE,
    paper_id UUID REFERENCES papers(id) ON DELETE CASCADE,
    PRIMARY KEY (review_id, paper_id)
);

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_papers_updated_at BEFORE UPDATE ON papers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trending_items_updated_at BEFORE UPDATE ON trending_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_learning_paths_updated_at BEFORE UPDATE ON learning_paths
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

