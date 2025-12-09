#!/bin/bash
set -e

# Seed CCF categories
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
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
EOSQL

echo "Seed data inserted successfully"

