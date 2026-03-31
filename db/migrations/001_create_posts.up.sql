create table if not exists posts
(
    id               varchar(36) primary key,
    title            varchar(200) not null,
    content          text         not null,
    author_id        varchar(36)  not null,
    comments_enabled boolean      not null default true,
    created_at       timestamptz  not null default current_timestamp,
    updated_at       timestamptz  not null default current_timestamp
);

create index if not exists idx_posts_author on posts (author_id);
create index if not exists idx_posts_created on posts (created_at, id);