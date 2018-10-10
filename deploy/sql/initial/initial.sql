
DROP TABLE IF EXISTS messenger_contacts;
CREATE TABLE messenger_contacts (
    id UUID DEFAULT gen_random_uuid(),
    member_id_1 UUID NOT NULL,
    member_id_2 UUID NOT NULL,
    state INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    UNIQUE INDEX uniq_friendship (member_id_1, member_id_2),
    INDEX idx_updated_at (updated_at)
);

DROP TABLE IF EXISTS messenger_groups;
CREATE TABLE messenger_groups (
    id UUID DEFAULT gen_random_uuid(),
    type INT NOT NULL,
    name STRING(24) NOT NULL,
    description STRING(256) NOT NULL,
    max_member_count INT NOT NULL,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    INDEX idx_updated_at (updated_at)
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
    id UUID DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    member_id UUID NOT NULL,    
    is_admin BOOL NOT NULL,   
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    UNIQUE INDEX uniq_member (group_id, member_id)
);

DROP TABLE IF EXISTS messenger_conversations;
CREATE TABLE messenger_conversations (
    id UUID DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    member_id UUID NOT NULL,
    is_mute BOOL NOT NULL,
    last_ack_message_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now()),
    PRIMARY KEY (id),
    UNIQUE INDEX uniq_member (group_id, member_id),
    INDEX idx_updated_at (updated_at)
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