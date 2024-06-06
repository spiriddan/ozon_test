CREATE TABLE IF NOT EXISTS post(
    id serial primary key,
    title text NOT NULL,
    body text,
    canComment boolean NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS comment(
    id serial primary key,
    body text NOT NULL,
    parent_post int REFERENCES post(id) ON DELETE CASCADE,
    parent_comment int REFERENCES comment(id) ON DELETE CASCADE,
    CHECK (
        (parent_post IS NOT NULL AND parent_comment IS NULL) OR
        (parent_post IS NULL AND parent_comment IS NOT NULL)
    )
)