-- tables

CREATE TABLE IF NOT EXISTS forums (
  slug         VARCHAR(255) NOT NULL PRIMARY KEY,
  title        VARCHAR(255) NOT NULL,
  "user"       VARCHAR(255) NOT NULL REFERENCES users (nickname),
  postsCount   INTEGER      NOT NULL DEFAULT 0,
  threadsCount INTEGER      NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS posts (
  id        BIGSERIAL    NOT NULL PRIMARY KEY,
  parent    VARCHAR(255) NOT NULL,
  author    VARCHAR(255) NOT NULL REFERENCES users (id),
  "message" TEXT         NOT NULL,
  isEdited  BOOLEAN      NOT NULL DEFAULT FALSE,
  forumSlug VARCHAR(255) NOT NULL REFERENCES forums (slug),
  threadId  INTEGER      NOT NULL REFERENCES threads (id),
  createdAt TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS threads (
  id         BIGSERIAL    NOT NULL PRIMARY KEY,
  title      VARCHAR(255) NOT NULL,
  author     VARCHAR(255) NOT NULL REFERENCES users (nickname),
  forumSlug  VARCHAR(255) NOT NULL,
  slug       VARCHAR(255) NOT NULL UNIQUE,
  "message"  TEXT         NOT NULL,
  votesCount INTEGER      NOT NULL DEFAULT 0,
  createdAt  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
  nickname VARCHAR(255) NOT NULL PRIMARY KEY,
  fullname VARCHAR(255) NOT NULL,
  email    VARCHAR(255) NOT NULL,
  about    TEXT         NOT NULL DEFAULT ''
);

-- indexes


-- triggers

CREATE OR REPLACE FUNCTION incrementForumsPostsCountFunc()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE forums
  SET postsCount = (postsCount + 1)
  WHERE slug = NEW.forumSlug;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER incrementForumsPostsCount
  AFTER INSERT OR DELETE
  ON posts
EXECUTE PROCEDURE incrementForumsPostsCountFunc();


CREATE OR REPLACE FUNCTION decrementForumsPostsCountFunc()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE forums
  SET postsCount = (postsCount - 1)
  WHERE slug = OLD.forumSlug;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER decrementForumsPostsCount
  AFTER INSERT OR DELETE
  ON posts
EXECUTE PROCEDURE decrementForumsPostsCountFunc();


CREATE OR REPLACE FUNCTION incrementForumsThreadsCountFunc()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE forums
  SET threadsCount = (threadsCount + 1)
  WHERE slug = NEW.forumSlug;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER incrementForumsThreadsCount
  AFTER INSERT OR DELETE
  ON posts
EXECUTE PROCEDURE incrementForumsThreadsCountFunc();


CREATE OR REPLACE FUNCTION decrementForumsThreadsCountFunc()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE forums
  SET threadsCount = (threadsCount - 1)
  WHERE slug = OLD.forumSlug;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER decrementForumsThreadsCount
  AFTER INSERT OR DELETE
  ON posts
EXECUTE PROCEDURE decrementForumsThreadsCountFunc();

