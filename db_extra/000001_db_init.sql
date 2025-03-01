CREATE TABLE users (
user_id    UUID    PRIMARY KEY NOT NULL,
first_name VARCHAR             NOT NULL,
last_name  VARCHAR             NOT NULL,
website    VARCHAR                 NULL,
file_name  VARCHAR                 NULL
);

CREATE TABLE user_creds (
user_id       UUID    PRIMARY KEY NOT NULL,
phone_number  VARCHAR                 NULL,
email         VARCHAR                 NULL,
password_hash VARCHAR             NOT NULL,
FOREIGN KEY(user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE posts (
post VARCHAR PRIMARY KEY
);

CREATE TABLE delivery_addresses (
user_id         UUID    PRIMARY KEY NOT NULL,
city            VARCHAR             NOT NULL,
post            VARCHAR             NOT NULL,
post_department INTEGER             NOT NULL,
FOREIGN KEY(user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY(post) REFERENCES posts(post) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE fundraise_statuses (
status VARCHAR PRIMARY KEY
);

CREATE TABLE fundraises (
fundraise_id  UUID    PRIMARY KEY      NOT NULL,
organizer_id  UUID                     NOT NULL,
title         VARCHAR                  NOT NULL,
description   VARCHAR                  NOT NULL,
target_amount NUMERIC(72, 18)          NOT NULL,
start_date    TIMESTAMP WITH TIME ZONE NOT NULL,
end_date      TIMESTAMP WITH TIME ZONE     NULL,
status        VARCHAR                  NOT NULL,
btc_address   VARCHAR                      NULL,
FOREIGN KEY(organizer_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY(status) REFERENCES fundraise_statuses(status) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE donations (
donation_id  UUID PRIMARY KEY         NOT NULL,
user_id      UUID                     NOT NULL,
fundraise_id UUID                     NOT NULL,
amount       NUMERIC(72, 18)          NOT NULL,
created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
FOREIGN KEY(user_id) REFERENCES users(user_id) ON UPDATE NO CASCADE ON DELETE NO ACTION,
FOREIGN KEY(fundraise_id) REFERENCES fundraises(fundraise_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE payments (
donation_id    UUID    PRIMARY KEY NOT NULL,
payment_type   VARCHAR             NOT NULL,
transaction_id VARCHAR             NOT NULL,
confirmed      BOOLEAN             NOT NULL,
FOREIGN KEY(donation_id) REFERENCES donations(donation_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE proofs (
fundraise_id UUID    PRIMARY KEY NOT NULL,
description  VARCHAR             NOT NULL,
FOREIGN KEY(fundraise_id) REFERENCES fundraises(fundraise_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE proof_images (
fundraise_id UUID    NOT NULL,
file_name    VARCHAR NOT NULL,
description  VARCHAR NOT NULL,
PRIMARY KEY(fundraise_id, file_name),
FOREIGN KEY(fundraise_id) REFERENCES proofs(fundraise_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE raffles (
raffle_id        UUID    PRIMARY KEY      NOT NULL,
title            VARCHAR                  NOT NULL,
description      VARCHAR                  NOT NULL,
minimum_donation NUMERIC(72, 18)          NOT NULL,
start_date       TIMESTAMP WITH TIME ZONE NOT NULL,
end_date         TIMESTAMP WITH TIME ZONE NOT NULL,
fundraise_id     UUID                     NOT NULL,
FOREIGN KEY(fundraise_id) REFERENCES fundraises(fundraise_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE gifts (
gift_id     UUID    PRIMARY KEY NOT NULL,
title       VARCHAR             NOT NULL,
description VARCHAR             NOT NULL,
raffle_id   UUID                NOT NULL,
user_id     UUID                NOT NULL,
FOREIGN KEY(raffle_id) REFERENCES raffles(raffle_id) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY(user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE NO ACTION
);

CREATE TABLE gift_images (
gift_id   UUID    NOT NULL,
file_name VARCHAR NOT NULL,
PRIMARY KEY(gift_id, file_name),
FOREIGN KEY(gift_id) REFERENCES gifts(gift_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE events (
event_id         UUID    PRIMARY KEY      NOT NULL,
title            VARCHAR                  NOT NULL,
description      VARCHAR                  NOT NULL,
date             DATE                     NOT NULL,
time             TIME                     NOT NULL,
format           VARCHAR                  NOT NULL,
max_participants INTEGER                  NOT NULL,
minimum_donation NUMERIC(72, 18)          NOT NULL,
address          VARCHAR                  NOT NULL,
status           VARCHAR                      NULL,
fundraise_id     UUID                     NOT NULL,
created_at       TIMESTAMP WITH TIME ZONE NOT NULL,
FOREIGN KEY(fundraise_id) REFERENCES fundraises(fundraise_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE event_participants (
event_id UUID NOT NULL,
user_id  UUID NOT NULL,
PRIMARY KEY(event_id, user_id),
FOREIGN KEY(event_id) REFERENCES events(event_id) ON UPDATE CASCADE ON DELETE CASCADE,
FOREIGN KEY(user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE NO ACTION
);

CREATE TABLE event_images (
event_id  UUID    NOT NULL,
file_name VARCHAR NOT NULL,
PRIMARY KEY(event_id, file_name),
FOREIGN KEY(event_id) REFERENCES events(event_id) ON UPDATE CASCADE ON DELETE CASCADE
);
