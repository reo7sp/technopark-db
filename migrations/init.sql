CREATE EXTENSION IF NOT EXISTS citext;

-- tables

CREATE TABLE IF NOT EXISTS users (
  nickname CITEXT       NOT NULL PRIMARY KEY,
  fullname VARCHAR(255) NOT NULL,
  email    CITEXT       NOT NULL UNIQUE,
  about    TEXT         NOT NULL DEFAULT ''
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
  id         BIGSERIAL NOT NULL PRIMARY KEY,
  parent     BIGINT REFERENCES posts (id),
  author     CITEXT    NOT NULL REFERENCES users (nickname),
  "message"  TEXT      NOT NULL,
  isEdited   BOOLEAN   NOT NULL DEFAULT FALSE,
  forumSlug  CITEXT    NOT NULL REFERENCES forums (slug),
  threadId   INTEGER   NOT NULL REFERENCES threads (id),
  threadSlug CITEXT REFERENCES threads (slug),
  createdAt  TIMESTAMP NOT NULL DEFAULT NOW(),
  path       BIGINT [] NOT NULL DEFAULT ARRAY [] :: BIGINT [],
  rootPostNo BIGINT
);

CREATE TABLE IF NOT EXISTS votes (
  nickname CITEXT   NOT NULL REFERENCES users (nickname),
  threadId BIGINT   NOT NULL REFERENCES threads (id),
  voice    SMALLINT NOT NULL,
  PRIMARY KEY (nickname, threadId)
);

-- indexes


-- triggers

CREATE OR REPLACE FUNCTION func_forums_incrementPostsCount()
  RETURNS TRIGGER AS
$$
BEGIN

  UPDATE forums
  SET postsCount = (postsCount + 1)
  WHERE slug = NEW.forumSlug;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_forums_incrementPostsCount
ON posts;

CREATE TRIGGER trig_forums_incrementPostsCount
  AFTER INSERT
  ON posts
  FOR EACH ROW
EXECUTE PROCEDURE func_forums_incrementPostsCount();


CREATE OR REPLACE FUNCTION func_forums_decrementPostsCount()
  RETURNS TRIGGER AS
$$
BEGIN

  UPDATE forums
  SET postsCount = (postsCount - 1)
  WHERE slug = OLD.forumSlug;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_forums_decrementPostsCount
ON posts;

CREATE TRIGGER trig_forums_decrementPostsCount
  AFTER DELETE
  ON posts
  FOR EACH ROW
EXECUTE PROCEDURE func_forums_decrementPostsCount();


CREATE OR REPLACE FUNCTION func_forums_incrementThreadsCount()
  RETURNS TRIGGER AS
$$
BEGIN

  UPDATE forums
  SET threadsCount = (threadsCount + 1)
  WHERE slug = NEW.forumSlug;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_forums_incrementThreadsCount
ON threads;

CREATE TRIGGER trig_forums_incrementThreadsCount
  AFTER INSERT
  ON threads
  FOR EACH ROW
EXECUTE PROCEDURE func_forums_incrementThreadsCount();


CREATE OR REPLACE FUNCTION func_forums_decrementThreadsCount()
  RETURNS TRIGGER AS
$$
BEGIN

  UPDATE forums
  SET threadsCount = (threadsCount - 1)
  WHERE slug = OLD.forumSlug;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_forums_decrementThreadsCount
ON threads;

CREATE TRIGGER trig_forums_decrementThreadsCount
  AFTER DELETE
  ON threads
  FOR EACH ROW
EXECUTE PROCEDURE func_forums_decrementThreadsCount();


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
    WHERE id = NEW.id;

    UPDATE posts
    SET path = ARRAY [NEW.id], rootPostNo = threadRootPostsCount
    WHERE id = NEW.id;

  ELSE

    SELECT rootPostsCount
    INTO threadRootPostsCount
    FROM threads
    WHERE threads.id = NEW.threadId;

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


CREATE OR REPLACE FUNCTION func_threads_insertVote()
  RETURNS TRIGGER
AS $$
BEGIN

  UPDATE threads
  SET votesCount = threads.votesCount + NEW.voice
  WHERE id = NEW.threadId;

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_threads_insertVote
ON votes;

CREATE TRIGGER trig_threads_insertVote
  AFTER INSERT
  ON votes
  FOR EACH ROW
EXECUTE PROCEDURE func_threads_insertVote();


CREATE OR REPLACE FUNCTION func_threads_updateVote()
  RETURNS TRIGGER
AS $$
BEGIN

  UPDATE threads
  SET votesCount = threads.votesCount + (NEW.voice - OLD.voice)
  WHERE id = NEW.threadId;

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_threads_updateVote
ON votes;

CREATE TRIGGER trig_threads_updateVote
  AFTER UPDATE
  ON votes
  FOR EACH ROW
EXECUTE PROCEDURE func_threads_updateVote();
