create table if not exists comments
(
    id         varchar(36) primary key,
    post_id    varchar(36)  not null references posts (id) on delete cascade,
    parent_id  varchar(36)  references comments (id) on delete cascade,
    content    text         not null,
    author_id  varchar(36)  not null,
    created_at timestamptz  not null default current_timestamp,
    updated_at timestamptz  not null default current_timestamp
);

create index if not exists idx_comments_post on comments (post_id);
create index if not exists idx_comments_parent on comments (parent_id);
create index if not exists idx_comments_created on comments (created_at, id);
create index if not exists idx_comments_post_parent on comments (post_id, parent_id);