CREATE EXTENSION IF NOT EXISTS citext;

-- tables

CREATE TABLE IF NOT EXISTS users (
  nickname CITEXT COLLATE ucs_basic     NOT NULL PRIMARY KEY,
  fullname VARCHAR(255)                 NOT NULL,
  email    CITEXT                       NOT NULL UNIQUE,
  about    TEXT                         NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS forums (
  slug         CITEXT       NOT NULL PRIMARY KEY,
  title        VARCHAR(255) NOT NULL,
  "user"       CITEXT       NOT NULL REFERENCES users (nickname),
  postsCount   INTEGER      NOT NULL DEFAULT 0,
  threadsCount INTEGER      NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads (
  id             BIGSERIAL                NOT NULL PRIMARY KEY,
  title          VARCHAR(255)             NOT NULL,
  author         CITEXT                   NOT NULL REFERENCES users (nickname),
  forumSlug      CITEXT                   NOT NULL REFERENCES forums (slug),
  slug           CITEXT UNIQUE,
  votesCount     INTEGER                  NOT NULL DEFAULT 0,
  "message"      TEXT                     NOT NULL,
  createdAt      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  rootPostsCount BIGINT                   NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts (
  id         BIGSERIAL                NOT NULL PRIMARY KEY,
  parent     BIGINT REFERENCES posts (id),
  author     CITEXT                   NOT NULL REFERENCES users (nickname),
  "message"  TEXT                     NOT NULL,
  isEdited   BOOLEAN                  NOT NULL DEFAULT FALSE,
  forumSlug  CITEXT                   NOT NULL REFERENCES forums (slug),
  threadId   INTEGER                  NOT NULL REFERENCES threads (id),
  threadSlug CITEXT REFERENCES threads (slug),
  createdAt  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  path       BIGINT []                NOT NULL DEFAULT ARRAY [] :: BIGINT [],
  rootPostNo BIGINT
);

CREATE TABLE IF NOT EXISTS votes (
  nickname CITEXT   NOT NULL REFERENCES users (nickname),
  threadId BIGINT   NOT NULL REFERENCES threads (id),
  voice    SMALLINT NOT NULL,
  PRIMARY KEY (nickname, threadId)
);

CREATE TABLE IF NOT EXISTS forumUsers (
  nickname  CITEXT NOT NULL REFERENCES users (nickname),
  forumSlug CITEXT NOT NULL REFERENCES forums (slug),
  PRIMARY KEY (nickname, forumSlug)
);

-- indexes

CREATE INDEX IF NOT EXISTS idx_users_nickname
  ON users (nickname);
CREATE INDEX IF NOT EXISTS idx_users_email
  ON users (email);

CREATE INDEX IF NOT EXISTS idx_forums_slug
  ON forums (slug);

CREATE INDEX IF NOT EXISTS idx_threads_id
  ON threads (id);
CREATE INDEX IF NOT EXISTS idx_threads_slug
  ON threads (slug);
CREATE INDEX IF NOT EXISTS idx_threads_forumSlugCreatedAt
  ON threads (forumSlug, createdAt);

CREATE INDEX IF NOT EXISTS idx_posts_id
  ON posts (id);
CREATE INDEX IF NOT EXISTS idx_posts_threadId
  ON posts (threadId);
CREATE INDEX IF NOT EXISTS idx_posts_threadSlug
  ON posts (threadSlug);
CREATE INDEX IF NOT EXISTS idx_posts_pathId
  ON posts (path, id);
CREATE INDEX IF NOT EXISTS idx_posts_rootPostNo
  ON posts (rootPostNo);
CREATE INDEX IF NOT EXISTS idx_posts_createdAtId
  ON posts (createdAt, id);

CREATE INDEX IF NOT EXISTS idx_votes_nicknameThreadId
  ON votes (nickname, threadId);

CREATE INDEX IF NOT EXISTS idx_forumUsers_forumSlugNickname
  ON forumUsers (forumSlug, nickname);

-- triggers

CREATE OR REPLACE FUNCTION func_posts_incrementForumPostsCount()
  RETURNS TRIGGER AS
$$
BEGIN

  UPDATE forums
  SET postsCount = (postsCount + 1)
  WHERE slug = NEW.forumSlug;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_posts_incrementForumPostsCount
ON posts;

CREATE TRIGGER trig_posts_incrementForumPostsCount
  AFTER INSERT
  ON posts
  FOR EACH ROW
EXECUTE PROCEDURE func_posts_incrementForumPostsCount();

-- CREATE OR REPLACE FUNCTION func_posts_decrementForumPostsCount()
--   RETURNS TRIGGER AS
-- $$
-- BEGIN
--
--   UPDATE forums
--   SET postsCount = (postsCount - 1)
--   WHERE slug = OLD.forumSlug;
--
--   RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- DROP TRIGGER IF EXISTS trig_posts_decrementForumPostsCount
-- ON posts;
--
-- CREATE TRIGGER trig_posts_decrementForumPostsCount
--   AFTER DELETE
--   ON posts
--   FOR EACH ROW
-- EXECUTE PROCEDURE func_posts_decrementForumPostsCount();


CREATE OR REPLACE FUNCTION func_threads_incrementForumThreadsCount()
  RETURNS TRIGGER AS
$$
BEGIN

  UPDATE forums
  SET threadsCount = (threadsCount + 1)
  WHERE slug = NEW.forumSlug;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_threads_incrementForumThreadsCount
ON threads;

CREATE TRIGGER trig_threads_incrementForumThreadsCount
  AFTER INSERT
  ON threads
  FOR EACH ROW
EXECUTE PROCEDURE func_threads_incrementForumThreadsCount();

-- CREATE OR REPLACE FUNCTION func_threads_decrementForumThreadsCount()
--   RETURNS TRIGGER AS
-- $$
-- BEGIN
--
--   UPDATE forums
--   SET threadsCount = (threadsCount - 1)
--   WHERE slug = OLD.forumSlug;
--
--   RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- DROP TRIGGER IF EXISTS trig_threads_decrementForumThreadsCount
-- ON threads;
--
-- CREATE TRIGGER trig_threads_decrementForumThreadsCount
--   AFTER DELETE
--   ON threads
--   FOR EACH ROW
-- EXECUTE PROCEDURE func_threads_decrementForumThreadsCount();


CREATE OR REPLACE FUNCTION func_posts_constructParentPath()
  RETURNS TRIGGER AS
$$
DECLARE
  threadRootPostsCount BIGINT;
  parentPath           BIGINT [];
BEGIN

  IF NEW.parent IS NULL
  THEN

    SELECT rootPostsCount
    INTO threadRootPostsCount
    FROM threads
    WHERE threads.id = NEW.threadId;

    UPDATE threads
    SET rootPostsCount = threadRootPostsCount + 1
    WHERE id = NEW.threadId;

    UPDATE posts
    SET path = ARRAY [NEW.id], rootPostNo = threadRootPostsCount
    WHERE id = NEW.id;

  ELSE

    SELECT rootPostNo
    INTO threadRootPostsCount
    FROM posts
    WHERE id = NEW.parent;

    SELECT path
    INTO parentPath
    FROM posts
    WHERE id = NEW.parent;

    UPDATE posts
    SET path = array_append(parentPath, NEW.id), rootPostNo = threadRootPostsCount
    WHERE id = NEW.id;

  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_posts_constructParentPath
ON posts;

CREATE TRIGGER trig_posts_constructParentPath
  AFTER INSERT
  ON posts
  FOR EACH ROW
EXECUTE PROCEDURE func_posts_constructParentPath();


CREATE OR REPLACE FUNCTION func_votes_insertThreadVote()
  RETURNS TRIGGER
AS $$
BEGIN

  UPDATE threads
  SET votesCount = threads.votesCount + NEW.voice
  WHERE id = NEW.threadId;

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_votes_insertThreadVote
ON votes;

CREATE TRIGGER trig_votes_insertThreadVote
  AFTER INSERT
  ON votes
  FOR EACH ROW
EXECUTE PROCEDURE func_votes_insertThreadVote();


CREATE OR REPLACE FUNCTION func_votes_updateThreadVote()
  RETURNS TRIGGER
AS $$
BEGIN

  UPDATE threads
  SET votesCount = threads.votesCount + (NEW.voice - OLD.voice)
  WHERE id = NEW.threadId;

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_votes_updateThreadVote
ON votes;

CREATE TRIGGER trig_votes_updateThreadVote
  AFTER UPDATE
  ON votes
  FOR EACH ROW
EXECUTE PROCEDURE func_votes_updateThreadVote();


CREATE OR REPLACE FUNCTION func_posts_insertForumUser()
  RETURNS TRIGGER
AS $$
BEGIN

  INSERT INTO forumUsers (nickname, forumSlug) VALUES (NEW.author, NEW.forumSlug)
  ON CONFLICT DO NOTHING;

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_posts_insertForumUser
ON posts;

CREATE TRIGGER trig_posts_insertForumUser
  AFTER INSERT
  ON posts
  FOR EACH ROW
EXECUTE PROCEDURE func_posts_insertForumUser();


CREATE OR REPLACE FUNCTION func_threads_insertForumUser()
  RETURNS TRIGGER
AS $$
BEGIN

  INSERT INTO forumUsers (nickname, forumSlug) VALUES (NEW.author, NEW.forumSlug)
  ON CONFLICT DO NOTHING;

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_threads_insertForumUser
ON threads;

CREATE TRIGGER trig_threads_insertForumUser
  AFTER INSERT
  ON threads
  FOR EACH ROW
EXECUTE PROCEDURE func_threads_insertForumUser();

-- settings

ALTER TABLE users
  SET ( AUTOVACUUM_VACUUM_THRESHOLD = 0),
  SET ( AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.0),
  SET ( AUTOVACUUM_FREEZE_MIN_AGE = 0),
  SET ( AUTOVACUUM_FREEZE_MAX_AGE = 100000),
  SET ( AUTOVACUUM_FREEZE_TABLE_AGE = 0);

ALTER TABLE forums
  SET ( AUTOVACUUM_VACUUM_THRESHOLD = 0),
  SET ( AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.0),
  SET ( AUTOVACUUM_FREEZE_MIN_AGE = 0),
  SET ( AUTOVACUUM_FREEZE_MAX_AGE = 100000),
  SET ( AUTOVACUUM_FREEZE_TABLE_AGE = 0);

ALTER TABLE threads
  SET ( AUTOVACUUM_VACUUM_THRESHOLD = 0),
  SET ( AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.0),
  SET ( AUTOVACUUM_FREEZE_MIN_AGE = 0),
  SET ( AUTOVACUUM_FREEZE_MAX_AGE = 100000),
  SET ( AUTOVACUUM_FREEZE_TABLE_AGE = 0);

ALTER TABLE posts
  SET ( AUTOVACUUM_VACUUM_THRESHOLD = 0),
  SET ( AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.0),
  SET ( AUTOVACUUM_FREEZE_MIN_AGE = 0),
  SET ( AUTOVACUUM_FREEZE_MAX_AGE = 100000),
  SET ( AUTOVACUUM_FREEZE_TABLE_AGE = 0);

ALTER TABLE votes
  SET ( AUTOVACUUM_VACUUM_THRESHOLD = 0),
  SET ( AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.0),
  SET ( AUTOVACUUM_FREEZE_MIN_AGE = 0),
  SET ( AUTOVACUUM_FREEZE_MAX_AGE = 100000),
  SET ( AUTOVACUUM_FREEZE_TABLE_AGE = 0);

ALTER TABLE forumUsers
  SET ( AUTOVACUUM_VACUUM_THRESHOLD = 0),
  SET ( AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.0),
  SET ( AUTOVACUUM_FREEZE_MIN_AGE = 0),
  SET ( AUTOVACUUM_FREEZE_MAX_AGE = 100000),
  SET ( AUTOVACUUM_FREEZE_TABLE_AGE = 0);
