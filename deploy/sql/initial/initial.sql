
DROP TABLE IF EXISTS messenger_contacts;
CREATE TABLE messenger_contacts (
    member_id UUID NOT NULL,
    friend_id UUID NOT NULL,
    name STRING(24) NOT NULL,
    state INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (member_id, friend_id)
);

DROP TABLE IF EXISTS messenger_groups;
CREATE TABLE messenger_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type INT NOT NULL,
    name STRING(24) NOT NULL,
    description STRING NOT NULL,
    max_member_count INT NOT NULL,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

DROP TABLE IF EXISTS messenger_groups_one;
CREATE TABLE messenger_groups_one (
    member_id_1 UUID NOT NULL,
    member_id_2 UUID NOT NULL,
    group_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (member_id_1, member_id_2)
);

DROP TABLE IF EXISTS messenger_group_members;
CREATE TABLE messenger_group_members (
    group_id UUID NOT NULL,
    member_id UUID NOT NULL,    
    is_admin BOOL NOT NULL,   
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (group_id, member_id)
);

DROP TABLE IF EXISTS messenger_conversations;
CREATE TABLE messenger_conversations (
    group_id UUID NOT NULL,
    member_id UUID NOT NULL,
    is_mute BOOL NOT NULL,
    last_ack_message_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (group_id, member_id)
);


DROP TABLE IF EXISTS messenger_messages;
CREATE TABLE messenger_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    type INT NOT NULL,
    state INT NOT NULL, 
    content STRING(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);