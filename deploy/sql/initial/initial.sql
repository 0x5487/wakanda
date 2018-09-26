DROP TABLE IF EXISTS friendship;
CREATE TABLE friendship (
    member_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    friend_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type INT,
    name STRING(24) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

DROP TABLE IF EXISTS group_members;
CREATE TABLE group_members (
    group_id UUID,
    member_id UUID,
    latest_ack_msgid INT,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);


DROP TABLE IF EXISTS group_messages;
CREATE TABLE group_messages (
    id INT PRIMARY KEY,
    group_id UUID,
    sender_id UUID,
    type INT,
    State INT, 
    name STRING(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);