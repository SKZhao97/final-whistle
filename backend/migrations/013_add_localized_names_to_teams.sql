-- Add localized Chinese display names for teams
ALTER TABLE teams
ADD COLUMN IF NOT EXISTS name_zh VARCHAR(100);

UPDATE teams
SET name_zh = CASE slug
  WHEN 'manchester-city' THEN '曼城'
  WHEN 'liverpool' THEN '利物浦'
  WHEN 'arsenal' THEN '阿森纳'
  WHEN 'chelsea' THEN '切尔西'
  WHEN 'manchester-united' THEN '曼联'
  WHEN 'tottenham-hotspur' THEN '热刺'
  ELSE name_zh
END
WHERE name_zh IS NULL OR name_zh = '';
