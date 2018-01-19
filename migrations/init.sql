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
  id         BIGSERIAL                NOT NULL PRIMARY KEY,
  title      VARCHAR(255)             NOT NULL,
  author     CITEXT                   NOT NULL REFERENCES users (nickname),
  forumSlug  CITEXT                   NOT NULL REFERENCES forums (slug),
  slug       CITEXT UNIQUE,
  votesCount INTEGER                  NOT NULL DEFAULT 0,
  "message"  TEXT                     NOT NULL,
  createdAt  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
  id         BIGSERIAL                NOT NULL PRIMARY KEY,
  parent     BIGINT REFERENCES posts (id),
  author     CITEXT                   NOT NULL REFERENCES users (nickname),
  "message"  TEXT                     NOT NULL,
  isEdited   BOOLEAN                  NOT NULL DEFAULT FALSE,
  forumSlug  CITEXT                   NOT NULL,
  threadId   INTEGER                  NOT NULL REFERENCES threads (id),
  createdAt  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  path       BIGINT [],
  rootParent BIGINT
);

CREATE TABLE IF NOT EXISTS votes (
  nickname CITEXT   NOT NULL REFERENCES users (nickname),
  threadId BIGINT   NOT NULL REFERENCES threads (id),
  voice    SMALLINT NOT NULL,
  PRIMARY KEY (nickname, threadId)
);

CREATE TABLE IF NOT EXISTS forumUsers (
  forumSlug CITEXT                   NOT NULL,
  nickname  CITEXT COLLATE ucs_basic NOT NULL,
  PRIMARY KEY (forumSlug, nickname)
);

-- indexes

CREATE UNIQUE INDEX IF NOT EXISTS idx_threads_slug_id
  ON threads (slug, id);
CREATE INDEX IF NOT EXISTS idx_threads_forumSlug_createdAt
  ON threads (forumSlug, createdAt);

CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_threadId_createdAt_id_asc
  ON posts (threadId, createdAt, id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_threadId_id_createdAt_id_asc
  ON posts (threadId, id, createdAt, id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_id_path
  ON posts (id, path);
CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_threadId_path_id
  ON posts (threadId, path, id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_threadId_path_id_parentnull
  ON posts (threadId, path, id)
  WHERE parent IS NULL;
CREATE INDEX IF NOT EXISTS idx_posts_rootParent
  ON posts (rootParent);

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


CREATE OR REPLACE FUNCTION func_posts_constructParentPath()
  RETURNS TRIGGER AS
$$
BEGIN

  IF NEW.parent IS NULL
  THEN

    UPDATE posts
    SET path = ARRAY [NEW.id], rootParent = NEW.id
    WHERE id = NEW.id;

  ELSE

    WITH parentPost AS (SELECT path
                        FROM posts
                        WHERE id = NEW.parent)
    UPDATE posts
    SET
      path       = array_append((SELECT path
                                 FROM parentPost), NEW.id),
      rootParent = (SELECT path [1]
                    FROM parentPost)
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

  INSERT INTO forumUsers (forumSlug, nickname) VALUES (NEW.forumSlug, NEW.author)
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

  INSERT INTO forumUsers (forumSlug, nickname) VALUES (NEW.forumSlug, NEW.author)
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
