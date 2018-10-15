
DROP TABLE IF EXISTS messenger_contacts;
CREATE TABLE messenger_contacts (
    group_id UUID NOT NULL,
    member_id_1 UUID NOT NULL,
    member_id_2 UUID NOT NULL,
    state INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (group_id),
    UNIQUE INDEX uniq_member (member_id_1, member_id_2),
    INDEX idx_member2 (member_id_2)
);

DROP TABLE IF EXISTS messenger_groups;
CREATE TABLE messenger_groups (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    type INT NOT NULL,
    name STRING(24) NOT NULL,
    description STRING(256) NOT NULL,
    max_member_count INT NOT NULL,
    member_count INT NOT NULL,
    state INT NOT NULL,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    INDEX idx_updated_at (updated_at)
);


DROP TABLE IF EXISTS messenger_group_members;
CREATE TABLE messenger_group_members (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    member_id UUID NOT NULL,    
    is_admin BOOL NOT NULL,   
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    UNIQUE INDEX uniq_member (member_id, group_id),
    INDEX idx_group (group_id)
);

DROP TABLE IF EXISTS messenger_conversations;
CREATE TABLE messenger_conversations (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    member_id UUID NOT NULL,
    is_mute BOOL NOT NULL,
    last_ack_message_id UUID NOT NULL,
    state INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    UNIQUE INDEX uniq_member (member_id, group_id)
);


DROP TABLE IF EXISTS messenger_messages;
CREATE TABLE messenger_messages (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    req_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    type INT NOT NULL,
    state INT NOT NULL, 
    content STRING(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    UNIQUE INDEX uniq_member (group_id, sender_id),
    UNIQUE INDEX uniq_request (req_id, sender_id),
    INDEX idx_updated_at (updated_at)
);