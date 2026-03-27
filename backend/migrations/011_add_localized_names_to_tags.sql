ALTER TABLE tags
ADD COLUMN IF NOT EXISTS name_en VARCHAR(50),
ADD COLUMN IF NOT EXISTS name_zh VARCHAR(50);

UPDATE tags
SET
  name_en = COALESCE(name_en, name),
  name_zh = COALESCE(
    name_zh,
    CASE slug
      WHEN 'hot-blooded' THEN '热血'
      WHEN 'boring' THEN '无聊'
      WHEN 'suffocating' THEN '窒息'
      WHEN 'classic' THEN '经典'
      WHEN 'unbelievable' THEN '离谱'
      WHEN 'regrettable' THEN '遗憾'
      WHEN 'dominance' THEN '统治'
      WHEN 'torture' THEN '折磨'
      WHEN 'comeback' THEN '逆转'
      WHEN 'destiny' THEN '宿命'
      ELSE name
    END
  );

ALTER TABLE tags
ALTER COLUMN name_en SET NOT NULL,
ALTER COLUMN name_zh SET NOT NULL;
