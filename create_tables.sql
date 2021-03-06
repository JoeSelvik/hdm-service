CREATE TABLE contenders (
  fb_id INTEGER,
  fb_group_id INTEGER,
  name TEXT,
  posts BLOB,
  avg_likes_per_post INTEGER,
  total_likes_received INTEGER,
  total_likes_given INTEGER,
  posts_used BLOB,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE posts (
  fb_id Text,
  fb_group_id INTEGER,
  posted_date DATETIME,
  author_fb_id INTEGER,
  likes BLOB,
  created_at DATETIME,
  updated_at DATETIME
);
